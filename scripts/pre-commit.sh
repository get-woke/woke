#!/usr/bin/env bash

set -e

if ! type woke > /dev/null; then
  echo "woke is not installed, or is not available in your PATH. See https://docs.getwoke.tech/installation."
  exit 1
fi

exec woke "${@}" --exit-1-on-failure
