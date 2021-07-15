#!/usr/bin/env bash
##
## This script installs all the tools needed for Go Development
##
(
    go install github.com/fzipp/gocyclo@latest
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v1.41.1
)
