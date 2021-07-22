<p align="center">
  <a href="https://github.com/get-woke/woke">
    <img alt="woke logo" src="https://raw.githubusercontent.com/get-woke/woke/main/assets/default-monochrome.svg" height="80" />
  </a>
  <h3 align="center">
    Detect non-inclusive language in your source code.
  </h3>
  <p align="center"><em>I stay woke - Erykah Badu</em></p>
</p>

---
<br />

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/get-woke/woke?logo=github&sort=semver)](https://github.com/get-woke/woke/releases)
![GitHub All Releases](https://img.shields.io/github/downloads/get-woke/woke/total)
[![Build](https://github.com/get-woke/woke/workflows/Build/badge.svg?branch=main)](https://github.com/get-woke/woke/actions)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/get-woke/woke)](https://goreportcard.com/report/github.com/get-woke/woke)
[![codecov](https://codecov.io/gh/get-woke/woke/branch/main/graph/badge.svg?token=BP133BM3HP)](https://codecov.io/gh/get-woke/woke/branch/main)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/get-woke/woke)](https://pkg.go.dev/github.com/get-woke/woke)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/get-woke/woke)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fget-woke%2Fwoke.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fget-woke%2Fwoke?ref=badge_shield)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#other-software)

---

Creating an inclusive work environment is imperative to a healthy, supportive, and
productive culture, and an environment where everyone feels welcome and included.

`woke` is a text file analysis tool that finds places within your source code that contain
non-inclusive language and suggests replacing them with more inclusive alternatives.

Companies like [GitHub](https://github.com/github/renaming), [Twitter](https://twitter.com/TwitterEng/status/1278733303508418560), and [Apple](https://developer.apple.com/news/?id=1o9zxsxl) are actively supporting a move to inclusive language.

![woke.gif](https://raw.githubusercontent.com/get-woke/get-woke.github.io/main/img/woke.gif)

**Table of Contents**

- [Installation](#installation)
  - [macOS](#macos)
  - [Windows](#windows)
  - [Simple installation](#simple-installation)
  - [Build from source](#build-from-source)
  - [Releases](#releases)
  - [Docker](#docker)
- [Usage](#usage)
  - [File globs](#file-globs)
  - [stdin](#stdin)
  - [Rules](#rules)
    - [Options](#options)
    - [Disabling Default Rules](#disabling-default-rules)
  - [Ignoring](#ignoring)
    - [Files](#files)
    - [`.wokeignore`](#wokeignore)
    - [In-line ignoring](#in-line-ignoring)
  - [Exit Code](#exit-code)
  - [Parallelism](#parallelism)
- [Tools](#tools)
- [Who uses `woke`](#who-uses-woke)
- [Resources](#resources)
- [Contributing](#contributing)
- [Versioning](#versioning)
- [Authors](#authors)
- [Acknowledgments](#acknowledgments)
- [License](#license)

## Installation

### macOS

You can install a binary release on macOS using [brew](https://brew.sh)

```bash
brew install get-woke/tap/woke
brew upgrade get-woke/tap/woke
```

### Windows

You can install `woke` with [`scoop`](https://scoop.sh/)

```sh
scoop bucket add get-woke https://github.com/get-woke/scoop-bucket.git
scoop install get-woke/woke
```

### Simple installation

To install the latest version:

```bash
curl -sSfL https://git.io/getwoke | \
  bash -s -- -b /usr/local/bin
```

Or install a specific version (omit the minor or patch portion to install the latest major/minor version)

```bash
curl -sSfL https://git.io/getwoke | \
  bash -s -- -b /usr/local/bin v0.9.0
```

Feel free to change the path from `/usr/local/bin`, just make sure `woke`
is available on your `$PATH` (check with `woke --version`).

### Build from source

Install the go toolchain: <https://golang.org/doc/install>

```bash
go get -u github.com/get-woke/woke
```

`woke` will be installed to `$GOPATH/bin/woke`.

### Releases

Download the latest binary from [Releases](https://github.com/get-woke/woke/releases/latest)

### Docker

You can run `woke` within docker. You will need to mount a volume that contains your source code and/or rules.

```bash
## Run with all defaults, within the mounted /src directory
docker run -v $(pwd):/src -w /src getwoke/woke

## Provide rules config
docker run -v $(pwd):/src -w /src getwoke/woke \
  woke -c my-rules.yaml
```

## Usage

```text
$ woke --help

woke is a linter that will check your source code for usage of non-inclusive
language and provide suggestions for alternatives. Rules can be customized
to suit your needs.

Provide a list file globs for files you'd like to check.

Usage:
  woke [globs ...] [flags]

Flags:
  -c, --config string       Config file (default is .woke.yaml in current directory, or $HOME)
      --debug               Enable debug logging
      --exit-1-on-failure   Exit with exit code 1 on failures
  -h, --help                help for woke
      --no-ignore           Ignored files in .gitignore/.wokeignore and inline ignores are processed
  -o, --output string       Output type [text,simple,github-actions,json] (default "text")
      --stdin               Read from stdin
  -v, --version             version for woke
```

### File globs

By default, `woke` will run against all text files in your current directory.
To change this, supply a space-separated list of globs as the first argument.

This can be something like `**/*.go`, or a space-separated list of filenames.

```bash
$ woke test.txt
test.txt:2:2-11: `Blacklist` may be insensitive, use `denylist`, `blocklist` instead (warning)
* Blacklist
  ^
test.txt:3:2-12: `White-list` may be insensitive, use `allowlist` instead (warning)
* White-list
  ^
test.txt:4:2-11: `whitelist` may be insensitive, use `allowlist` instead (warning)
* whitelist
  ^
test.txt:5:2-11: `blacklist` may be insensitive, use `denylist`, `blocklist` instead (warning)
* blacklist
  ^
```

### stdin

You can also provide text to `woke` via stdin

```bash
$ echo "This has whitelist from stdin" | woke --stdin
/dev/stdin:1:9-18: `whitelist` may be insensitive, use `allowlist` instead (warning)
This has whitelist from stdin
         ^
```

### Rules

Configure your custom rules config in `.woke.yaml` or `.woke.yml`. `woke` uses the following precedence order. Each item takes precedence over the item below it:

- `current working directory`
- `$HOME`

This file will be picked up automatically up your customizations without needing to supply it with the `-c` flag.

See [example.yaml](https://github.com/get-woke/woke/blob/main/example.yaml) for an example of adding custom rules.
You can also supply your own rules with `-c path/to/rules.yaml` if you want to handle different rulesets.

The syntax for rules is very basic. You just need a name, a list of terms to match that violate the rule,
and a list of alternative suggestions.

```yaml
rules:
  - name: whitelist
    terms:
      - whitelist
      - white-list
    alternatives:
      - allowlist
    note: An optional description why these terms are not inclusive. It can be optionally included in the output message.
    # options:
    #   word_boundary: false
    #   include_note: false
```

A set of default rules is provided in [`pkg/rule/default.yaml`](https://github.com/get-woke/woke/blob/main/pkg/rule/default.yaml).

**NOTE: if you copy these rules into your config file, be sure to put them under the `rules:` key.**

#### Options

You can configure options for each rule. Add an `options` key to your rule definition to customize.

Current options include:

- `word_boundary` (default: `false`)
  - If `true`, terms will trigger findings when they are surrounded by ASCII word boundaries.
  - If `false`, will trigger findings if the term if found anywhere in the line, regardless if it is within an ASCII word boundary.
- `include_note` (default: `not set`)
  - If `true`, the rule note will be included in the output message explaining why this finding is not inclusive
  - If `false`, the rule note will not be included in the output message
  - If `not set`, `include_note` in your `woke` config file (ie `.woke.yml`) regulates if the note should be included in the output message (default: `false`).

#### Disabling Default Rules

You can disable default rules by providing a rule in your `woke` config file (ie `.woke.yml`), with no terms or alternatives.

This will disable the default `whitelist` rule:

```yaml
rules:
  - name: whitelist
```

### Ignoring

#### Files

You can ignore files by adding to your config file. All ways of ignoring files below should follow the [gitignore](https://git-scm.com/docs/gitignore) convention.

```yaml
ignore_files:
  - .git/*
  - other/files/in/repo
  - globs/too/*
```

`woke` will also automatically ignore anything listed in `.gitignore`.

#### `.wokeignore`

You may also specify a `.wokeignore` file at the root of the directory to add additional ignore files.
This also follows the [gitignore](https://git-scm.com/docs/gitignore) convention.

See [.wokeignore.example](.wokeignore.example) for a collection of common files and directories that may contain generated [SHA](https://en.wikipedia.org/wiki/Secure_Hash_Algorithms) and [GUID](https://en.wikipedia.org/wiki/Universally_unique_identifier)s.

#### In-line ignoring

There may be times where you don't want to ignore an entire file.
You may ignore a specific line for one or more rules by creating an in-line comment.

This functionality is very rudimentary, it does a simple search for the phrase. Since
`woke` is just a text file analyzer, it has no concept of the comment syntax for every file
type it might encounter.

Simply add the following to the line you wish to ignore, using comment syntax that is supported for your file type.
(`woke` is not responsible for broken code due to in-line ignoring. Make sure you comment correctly!)

```bash
This line has RULE_NAME but will be ignored # wokeignore:rule=RULE_NAME

# for example, to ignore the following line for the whitelist rule
whitelist # wokeignore:rule=whitelist

# or for multiple rules
whitelist and blacklist # wokeignore:rule=whitelist,blacklist
```

Here's an example in go:

```go
func main() {
  fmt.Println("here is the whitelist") // wokeignore:rule=whitelist
}
```

### Exit Code

By default, `woke` will exit with a successful exit code when there are any rule failures.
The idea is, if you run `woke` on PRs, you may not want to block a merge, but you do
want to inform the author that they can make better word choices.

If you're using `woke` on PRs, you can choose to enforce these rules with a non-zero
exit code by running `woke --exit-1-on-failure`.

### Parallelism

By default, `woke` will parse files in parallel and will consume as many resources as it can, unbounded.
This means `woke` will be fast, but might run out of memory, depending on how large the files/lines are.

We can limit these allocations by bounding the number of files read in parallel. To accomplish this,
set the environment variable `WORKER_POOL_COUNT` to an integer value of the fixed number of goroutines
you would like to spawn for reading files.

Read more about go's concurrency patterns [here](https://blog.golang.org/pipelines).

## Tools

- [GitHub Action](https://github.com/marketplace/actions/run-woke)
- [GitHub Action (reviewdog)](https://github.com/marketplace/actions/run-woke-with-reviewdog)
- [VS Code Extension](https://marketplace.visualstudio.com/items?itemName=get-woke.vscode-woke)

## Who uses `woke`

Are you using `woke` in your org? We'd love to know! Please send a PR and add your organization name/link to this list.

- [Rent the Runway](https://www.renttherunway.com) ([blog post](https://dresscode.renttherunway.com/blog/woke))
- [Veue](https://veuelive.com)

## Resources

- <https://buffer.com/resources/inclusive-language-tech/>
- <https://medium.com/pm101/inclusive-language-guide-for-tech-companies-and-startups-f5b254d4a5b7>
- <https://www.marketplace.org/2020/06/17/tech-companies-update-language-to-avoid-offensive-terms/>
- <https://tools.ietf.org/html/draft-knodel-terminology-02>

## Contributing

Please read [CONTRIBUTING.md](https://github.com/get-woke/woke/blob/main/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/get-woke/woke/tags).

## Authors

- **Caitlin Elfring** - [caitlinelfring](https://github.com/caitlinelfring)

<!-- See also the list of [contributors](https://github.com/get-woke/woke/contributors) who participated in this project. -->

## Acknowledgments

The following projects provided inspiration for parts of `woke`

- <https://github.com/get-alex/alex>
- <https://github.com/retextjs/retext-equality>
- <https://github.com/golangci/golangci-lint>

## License

This application is licensed under the MIT License, you may obtain a copy of it
[here](https://github.com/get-woke/woke/blob/main/LICENSE).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fget-woke%2Fwoke.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fget-woke%2Fwoke?ref=badge_large)
