package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	openai "github.com/sashabaranov/go-openai"
)

const (
	systemPrompt = `
	You are a highly skilled systems engineer with expertise in crafting flawless shell scripts. 
	Your primary responsibility involves two key tasks: Firstly, receiving a shell script that has 
	been evaluated by a static analysis tool; and secondly, analyzing the provided list of warnings 
	and issues identified by this tool. Utilizing these inputs, your objective is to develop an 
	enhanced version of the shell script that effectively addresses and resolves all identified 
	warnings and issues.
	
	Your task is to revise a shell script strictly based on the warnings and errors highlighted by 
	the static analysis tool. You are required to make changes only to address these specific issues. 
	Any modifications that do not directly or indirectly relate to correcting these warnings and errors 
	should be avoided. Additionally, you must preserve the original comments in the script, altering 
	them only if they become irrelevant due to the changes you implement.

	Your response should exclusively consist of the updated shell script text, presented without using 
	a code block format.
	`
	
	userPrompt = `
	SHELL_SCRIPT:
	{{.ScriptContents}}
	
	STATIC_ANALYSIS_OUTPUT:
	{{.StaticAnalysisOutput}}
	
	Your task is to revise the provided shell script, focusing solely on rectifying the warnings and errors 
	identified in the static analysis output. 
	Ensure that the output of your task is solely the modified shell script text, presented without the use 
	of a code block format.
	`

	gpt35turbo = "gpt-3.5-turbo"
	gpt4turbo = "gpt-4-turbo"
)

var (
	userPromptTmpl = template.Must(template.New("prompt").Parse(userPrompt))
	client = openai.NewClient(os.Getenv("OPENAI_API_KEY")) 
)

var (
	// Command line flags
	writeFile bool
	showVersion bool
	model string

	version = "dev"
)

func init() {
	flag.BoolVar(&writeFile, "w", false, "write shell script to input file")
	flag.BoolVar(&showVersion, "v", false, "print version number and exit")
	flag.StringVar(&model, "m", gpt35turbo, fmt.Sprintf("specify the model to use (%s or %s)", gpt35turbo, gpt4turbo))

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] FILE\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Execute shellcheck on the given script and pass the results to a large language model " +
		"for making appropriate corrections.\n\n")
	fmt.Fprintf(os.Stderr, "The default behavior displays the modified script on the console. Use the '-w' flag " +
		"to save the changes directly to the specified file.\n\n")
		fmt.Fprintf(os.Stderr, "The shellcheck binary must be present in your path.\n\n")
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "ENVIRONMENT:")
	fmt.Fprintln(os.Stderr, "  OPENAI_API_KEY OpenAI API key")
}

func printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(color.Output, format, a...)
}

func main() {
	flag.Parse()
	
	if showVersion {
		fmt.Printf("%s %s (runtime: %s)\n", os.Args[0], version, runtime.Version())
		os.Exit(0)
	}

	if model != gpt35turbo && model != gpt4turbo {
		fmt.Fprintf(os.Stderr, "%s: model must be %s or %s\n", os.Args[0], gpt35turbo, gpt4turbo)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filePath := args[0]

	run(filePath)
}

func run(filePath string) {
	printf("%s\n", "Running shellcheck against script")
	analysis, err := runShellcheck(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if analysis == "" {
		printf("%s\n", color.GreenString("shellcheck detected no issues."))
		return
	}

	printf("%s\n", color.YellowString("shellcheck detected potential issues"))

	script, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	result, err := callCompletionAPI(string(script), string(analysis))
	if err != nil {
		log.Fatalf("error calling completion API: %v", err)
	}

	if writeFile {
		if err := os.WriteFile(filePath, []byte(result), 0644); err != nil {
			log.Fatalf("could not write updated script to file: %v", err)
		}
		printf("%s %s\n", color.GreenString("Updated script written to"), color.GreenString(filePath))
		printf("%s\n", color.YellowString("Double check it before you commit!"))
	} else {
		printf("%s\n", result)
	}
}

func runShellcheck(filePath string) (string, error) {
	cmd := exec.Command("shellcheck", filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		// shellcheck returns exit code 1 when it finds issues
		if exitCode == 1 {
			return string(output), nil
		} else {	
			return "", fmt.Errorf("shellcheck exited with code %d: %s", exitCode, err)
		}
	}

	return "", nil
}

func callCompletionAPI(script, analysis string) (string, error) {
	spin := spinner.New(spinner.CharSets[26], 250*time.Millisecond)
	spin.Prefix = "Waiting for completion response"
	spin.Start()
	defer spin.Stop()

	data := map[string]string{
		"ScriptContents":       string(script),
		"StaticAnalysisOutput": string(analysis),
	}

	var buffer bytes.Buffer
	err := userPromptTmpl.Execute(&buffer, data)
	if err != nil {
		log.Fatalf("unable to format prompt: %v", err)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		getCompletionRequest(buffer.String()),
	)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("empty response")
}

func getCompletionRequest(prompt string) openai.ChatCompletionRequest {
	m := openai.GPT3Dot5Turbo
	if model == gpt4turbo {
		m = openai.GPT4TurboPreview
	}

	return openai.ChatCompletionRequest{
		Model: m,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
} 