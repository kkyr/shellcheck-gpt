# shellcheck-gpt

Automatically correct script issues by feeding [shellcheck](https://www.shellcheck.net) analysis into an LLM.

## Getting started

### Prerequisites

- [ShellCheck](https://www.shellcheck.net) should be installed and in your $PATH.
- [Go](https://go.dev/doc/install)
- A valid OpenAI API key

### Basic installation

shellcheck-gpt can be installed using the `go install` tool. Distribution methods that do not depend on Go are coming soon!

```shell
$ go install github.com/kkyr/shellcheck-gpt
```

## Usage

Store your OpenAI API key into the environment:

```shell
export OPENAI_API_KEY=replace-with-your-key
```

Run shellcheck-gpt against a script:

```shell
$ shellcheck-gpt script.sh
```

This will:

1. Run shellcheck against `script.sh`
1. Feed the script and the output of shellcheck into an LLM and ask it to make the corrections
1. Write the LLM's output onto stdout

If you'd like to write the output back into the script, use the `-w` flag:

```shell
$ shellcheck-gpt -w script.sh
```

## Configuration

By default, shellcheck-gpt uses the gpt-3.5-turbo model. You can modify this to another model using the `-m` flag:

```shell
$ shellcheck-gpt -m gpt-4-turbo script.sh
```

See available options and models using the `--help` flag:

```shell
$ shellcheck-gpt --help
```

## Contributing

Contributions are welcome!
