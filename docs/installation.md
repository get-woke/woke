# Installation

!!! tip
    There are multiple ways to install `woke`. If you're interested in any installation methods
    that are not listed here, feel free to open an [issue]({{config.repo_url}}issues).

## Releases

Download the latest binary from [Releases]({{config.repo_url}}releases/latest)

## macOS

You can install a binary release on macOS using [brew](https://brew.sh)

```bash
brew install get-woke/tap/woke
brew upgrade get-woke/tap/woke
```

## Windows

You can install `woke` with [`scoop`](https://scoop.sh/)

```sh
scoop bucket add get-woke https://github.com/get-woke/scoop-bucket.git
scoop install get-woke/woke
```

## Simple installation

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

## Build from source

Install the go toolchain: <https://golang.org/doc/install>

```bash
go install github.com/get-woke/woke@latest

# Or pin a specific version
go install github.com/get-woke/woke@v0.9.0
```

## Docker

You can run `woke` within docker. You will need to mount a volume that contains your source code and/or rules.

```bash
## Run with all defaults, within the mounted /src directory
docker run -v $(pwd):/src -w /src getwoke/woke

## Provide rules config
docker run -v $(pwd):/src -w /src getwoke/woke \
  woke -c my-rules.yaml
```

## CI

### GitHub Actions

- [GitHub Action](https://github.com/marketplace/actions/run-woke)
- [GitHub Action (reviewdog)](https://github.com/marketplace/actions/run-woke-with-reviewdog)

### Others

Are there other CI systems that you're using to run `woke`? Edit this page and add documentation/configurations for others.
