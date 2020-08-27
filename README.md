# woke

_I stay woke - Erykah Badu_

[![Build](https://github.com/caitlinelfring/woke/workflows/Build/badge.svg?branch=main)](https://github.com/caitlinelfring/woke/actions)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/caitlinelfring/woke)](https://pkg.go.dev/github.com/caitlinelfring/woke)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/caitlinelfring/woke)](https://goreportcard.com/report/github.com/caitlinelfring/woke)

`woke` is a text file analysis tool that detects non-inclusive language in your source code.

![woke.gif](https://raw.githubusercontent.com/caitlinelfring/woke/main/img/woke.gif)

## About

Creating an inclusive work environment is imperitive to a healthy, supportive, and
productive culture, and an environment where everyone feels welcome and included.

`woke`'s purpose is to point out places where improvements can be made by removing
 non-inclusive language and replacing it with more inclusive alternatives.

Companies like [GitHub](https://github.com/github/renaming), [Twitter](https://twitter.com/TwitterEng/status/1278733303508418560), and [Apple](https://developer.apple.com/news/?id=1o9zxsxl) are already pushing these changes.

## Installation

```bash
go get -u github.com/caitlinelfring/woke
```

`woke` will be installed to `$GOPATH/bin/woke`.

Alternatively, download the latest binary from [Releases](https://github.com/caitlinelfring/woke/releases/latest)

### Docker

You can run `woke` within docker. You will need to mount a volume that contains your source code and/or rules.

```bash
## Run with all defaults, within the mounted /src directory
docker run -v $(pwd):/src -w /src celfring/woke

## Provide rules config
docker run -v $(pwd):/src -w /src celfring/woke \
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
  -o, --output string       Output type [text,simple] (default "text")
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
        4:2-4:11       warn        Instead of 'whitelist', consider the following alternative(s): 'allowlist'
        5:2-5:11       warn        Instead of 'blacklist', consider the following alternative(s): 'denylist,blocklist'
```

### stdin

You can also provide text to `woke` via stdin

```bash
$ echo "This has whitelist from stdin" | woke --stdin
/dev/stdin
        1:8-1:17       warn       Instead of 'whitelist', consider the following alternative(s): 'allowlist'
```

### Rules

A set of default rules is provided in [`example.yaml`](https://github.com/caitlinelfring/woke/blob/main/example.yaml).
You can supply your own rules with `-c path/to/rules.yaml`

The syntax for rules is very basic. You just need a name, a regex used
to match, and a string of alternatives.

```yaml
rules:
  - name: whitelist
    regexp: \b(white-?list)\b
    alternatives: allowlist
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

If you're using `work` on PRs, you can choose to enforce these rules with a non-zero
exit code, but running `work --exit-1-on-failure`.

## TODO

* Benchmarking
* Tests
* What happens when run on a large repo?
* GitHub Actions
* More rules

## Resources

* <https://buffer.com/resources/inclusive-language-tech/>
* <https://medium.com/pm101/inclusive-language-guide-for-tech-companies-and-startups-f5b254d4a5b7>
* <https://www.marketplace.org/2020/06/17/tech-companies-update-language-to-avoid-offensive-terms/>
* <https://tools.ietf.org/html/draft-knodel-terminology-02>

## License

This application is licensed under the MIT License, you may obtain a copy of it
[here](https://github.com/caitlinelfring/woke/blob/main/LICENSE).
