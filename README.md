# woke

_I stay woke - Erykah Badu_

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/get-woke/woke?logo=github&sort=semver)](https://github.com/get-woke/woke/releases)
[![Build](https://github.com/get-woke/woke/workflows/Build/badge.svg?branch=main)](https://github.com/get-woke/woke/actions)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/get-woke/woke)](https://pkg.go.dev/github.com/get-woke/woke)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/get-woke/woke)](https://goreportcard.com/report/github.com/get-woke/woke)
[![Coverage Status](https://coveralls.io/repos/github/get-woke/woke/badge.svg?branch=main)](https://coveralls.io/github/get-woke/woke?branch=main)

`woke` is a text file analysis tool that detects non-inclusive language in your source code.

**Table of Contents**

- [woke](#woke)
  - [About](#about)
  - [Installation](#installation)
    - [Simple installation](#simple-installation)
    - [Build from source](#build-from-source)
    - [Releases](#releases)
    - [Docker](#docker)
  - [Usage](#usage)
    - [File globs](#file-globs)
    - [stdin](#stdin)
    - [Rules](#rules)
    - [Ignoring files](#ignoring-files)
      - [`.wokeignore`](#wokeignore)
    - [Exit Code](#exit-code)
  - [Tools](#tools)
  - [TODO](#todo)
  - [Resources](#resources)
  - [License](#license)

## About

Creating an inclusive work environment is imperitive to a healthy, supportive, and
productive culture, and an environment where everyone feels welcome and included.

`woke`'s purpose is to point out places where improvements can be made by removing
 non-inclusive language and replacing it with more inclusive alternatives.

Companies like [GitHub](https://github.com/github/renaming), [Twitter](https://twitter.com/TwitterEng/status/1278733303508418560), and [Apple](https://developer.apple.com/news/?id=1o9zxsxl) are already pushing these changes.

## Installation

### Simple installation

```bash
curl -sSfL https://git.io/getwoke | \
  bash -s -- -b /usr/local/bin
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
  -c, --config string       YAML file with list of rules
      --debug               Enable debug logging
      --exit-1-on-failure   Exit with exit code 1 on failures
  -h, --help                help for woke
  -o, --output string       Output type [text,simple,github-actions] (default "text")
      --stdin               Read from stdin
  -v, --version             version for woke
```

### File globs

By default, `woke` will run against all text files in your current directory.
To change this, supply a space-separated list of globs as the first argument.

This can be something like `**/*.go`, or a space-separated list of filenames.

```bash
$ woke test.txt
test.txt
        2:2-2:11       warn        `Blacklist` may be insensitive, use `blocklist` instead
        3:2-3:12       warn        `White-list` may be insensitive, use `allowlist` instead
        4:2-4:11       warn        `whitelist` may be insensitive, use `allowlist` instead
        5:2-5:11       warn        `blacklist` may be insensitive, use `blocklist` instead
```

### stdin

You can also provide text to `woke` via stdin

```bash
$ echo "This has whitelist from stdin" | woke --stdin
/dev/stdin
        1:9-1:18       warn        `whitelist` may be insensitive, use `allowlist` instead
```

### Rules

A set of default rules is provided in [`pkg/rule/default.go`](https://github.com/get-woke/woke/blob/main/pkg/rule/default.go).

See [example.yaml](https://github.com/get-woke/woke/blob/example.yaml) for an example of adding custom rules.
You can supply your own rules with `-c path/to/rules.yaml`

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
```

### Ignoring files

In your config file, you can ignore files by adding:

```yaml
ignore_files:
  - .git/*
  - other/files/in/repo
  - globs/too/*
```

`woke` will also automatically ignore anything listed in `.gitignore`.

`woke` uses [go-gitignore](github.com/sabhiram/go-gitignore) to ignores.
This follows the common `.gitignore` convention. See link for more details on matching.

#### `.wokeignore`

You may also specify a `.wokeignore` file at the root of the directory to add additional ignore files.
This also follows the `.gitignore` convention.

### Exit Code

By default, `woke` will exit with a successful exit code when there are any rule failures.
The idea is, if you run `woke` on PRs, you may not want to block a merge, but you do
want to inform the author that they can make better word choices.

If you're using `woke` on PRs, you can choose to enforce these rules with a non-zero
exit code, but running `woke --exit-1-on-failure`.

## Tools

- [GitHub Action](https://github.com/marketplace/actions/run-woke)
- [GitHub Action (reviewdog)](https://github.com/marketplace/actions/run-woke-with-reviewdog)

## TODO

* Benchmarking
  * What happens when run on a large repo?
* More rules

## Resources

* <https://buffer.com/resources/inclusive-language-tech/>
* <https://medium.com/pm101/inclusive-language-guide-for-tech-companies-and-startups-f5b254d4a5b7>
* <https://www.marketplace.org/2020/06/17/tech-companies-update-language-to-avoid-offensive-terms/>
* <https://tools.ietf.org/html/draft-knodel-terminology-02>

## License

This application is licensed under the MIT License, you may obtain a copy of it
[here](https://github.com/get-woke/woke/blob/main/LICENSE).
