# woke

[![GoDoc](https://godoc.org/github.com/caitlinelfring/woke?status.svg)](https://godoc.org/github.com/caitlinelfring/woke)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](http://choosealicense.com/licenses/mit/)

`woke` is a text file analysis tool that detects non-inclusive language in your source code.

Creating an inclusive work environment is imperitive to a healthy, supportive, and
productive culture, and an environment where everyone feels welcome and included.

`woke`'s purpose is to point out places where improvements can be made by removing
 non-inclusive language and replacing it with more inclusive alternatives.

Companies like [GitHub](https://github.com/github/renaming), [Twitter](https://twitter.com/TwitterEng/status/1278733303508418560), and [Apple](https://developer.apple.com/news/?id=1o9zxsxl) are already pushing these changes.

## Installation

Using `woke` is straightforward. First, use `go get` to install the latest version
_(requires go, more instructions for non-go installation coming soon)_.

```bash
go get -u github.com/caitlinelfring/woke
```

`woke` should be installed to `$GOPATH/bin/woke`.

## Usage

```bash
$ woke --help

woke is a linter that will check your source code for usage of non-inclusive
language and provide suggestions for alternatives. Rules can be customized
to suit your needs.

Provide a list of comma-separated file globs for files you'd like to check.

Usage:
  woke (file globs to check) [flags]

Flags:
      --exit-1-on-failure    Exit with exit code 1 on failures. Otherwise, will always exit 0 if any failures occur
  -h, --help                 help for woke
  -r, --rule-config string   YAML file with list of rules (default "default.yaml")
```

### File globs

By default, `woke` will run against all text files in your current directory.
To change this, supply a comma-separated list of globs as the first argument.

This can be something like `**/*.go`, or a comma-separated list of filenames.

```bash
$ woke test.txt
[test.txt:2:2] Instead of 'Blacklist', consider the following alternative(s): 'denylist,blocklist'
[test.txt:3:2] Instead of 'White-list', consider the following alternative(s): 'allowlist'
[test.txt:4:2] Instead of 'whitelist', consider the following alternative(s): 'allowlist'
[test.txt:5:2] Instead of 'blacklist', consider the following alternative(s): 'denylist,blocklist'
```

### Rules

A set of default rules is provided in [`example.yaml`](example.yaml).
You can supply your own rules with `-r path/to/rules.yaml`

The syntax for rules is very basic. You just need a name, a regex used
to match, and a string of alternatives.

```yaml
rules:
  - name: whitelist
    regexp: \b(white-?list)\b
    alternatives: allowlist
```

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
* Color output
* Warn/error rule severities

## Resources

* <https://buffer.com/resources/inclusive-language-tech/>
* <https://medium.com/pm101/inclusive-language-guide-for-tech-companies-and-startups-f5b254d4a5b7>
* <https://www.marketplace.org/2020/06/17/tech-companies-update-language-to-avoid-offensive-terms/>
* <https://tools.ietf.org/html/draft-knodel-terminology-02>

## License

This application is licensed under the MIT License, you may obtain a copy of it [here](LICENSE).
