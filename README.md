<p align="center">
  <img src="resources/logo.png" alt="shellcheck-gpt" title="shellcheck-gpt" class="img-responsive" style="width:256px;" />
  <h3 align="center">shellcheck-gpt</h3>
  <p align="center">Instantly enhance nad optimize shell scripts with the power of ShellCheck and LLMs.</p>
</p>

<p align="center">
  <a href="https://github.com/kkyr/shellcheck-gpt/releases"><img src="https://img.shields.io/github/v/tag/kkyr/shellcheck-gpt?style=for-the-badge" alt="semver tag" title="semver tag"/></a>
  <a href="https://github.com/kkyr/shellcheck-gpt/actions?workflow=build"><img src="https://img.shields.io/github/actions/workflow/status/kkyr/shellcheck-gpt/build.yml?style=for-the-badge&branch=main" alt="build status" title="semver tag"/></a>
  <a href="https://github.com/kkyr/shellcheck-gpt/blob/master/LICENSE"><img src="https://img.shields.io/github/license/kkyr/shellcheck-gpt?style=for-the-badge" alt="license" title="license"/></a>
</p>

---

# Demo

![Alt Text](https://github.com/kkyr/shellcheck-gpt/blob/main/resources/screencast.gif)

_Access screencast directly [here](https://asciinema.org/a/fhFJ4gOnTmscpXiv1msPfxyAz)_.

See [example](https://github.com/kkyr/shellcheck-gpt/tree/main/example) directory to see script contents before and after.

# Getting started

## Prerequisites

- [ShellCheck](https://www.shellcheck.net) should be installed and in your $PATH.
- An active OpenAI API key.

## Installation

### Homebrew

```shell
brew install kkyr/tap/shellcheck-gpt
```

### Pre-built binaries

Download the latest [release](https://github.com/kkyr/shellcheck-gpt/releases) and add the executable to your PATH.

### Build using Go toolchain

```shell
go install github.com/kkyr/shellcheck-gpt
```

# Usage

Add your OpenAI API key to the environment:

```shell
export OPENAI_API_KEY=replace-with-your-api-key
```

Run shellcheck-gpt against a script:

```shell
shellcheck-gpt script.sh
```

This will:

1. Run shellcheck against `script.sh`
1. Feed the contents of the script and the output of shellcheck into an OpenAI LLM and ask it to make corrections
1. Write the LLM's output to stdout

> [!WARNING]  
> The entire content of your script is sent in cleartext to OpenAI.

If you'd like to instead write the output back into the script, use the `-w` flag:

```shell
shellcheck-gpt -w script.sh
```

# Configuration

By default, shellcheck-gpt uses the _gpt-3.5-turbo_ model. You change model using the `-m` flag:

```shell
shellcheck-gpt -m gpt-4-turbo script.sh
```

See available options and models using the `--help` flag:

```shell
shellcheck-gpt --help
```

# Contributing

Contributions are welcome!

Areas for improvement:

- [ ] Use `shellcheck -f diff -p1` to apply fixes that can be automatically applied before calling into LLM.
- [ ] Support more LLMs
- [ ] Add verbose flag to print LLM input/output
- [ ] Refactor and encapsulate components
