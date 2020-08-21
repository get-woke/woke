#!/bin/bash
# Autoreloads go application on file changes
# From https://github.com/caitlinelfring/go_autoreload

TARGET=${1:-"app"}
FILES_TO_WATCH=${2:-"./*.go"}

wait_for_changes() {
  echo "Waiting for changes"
  if [ "$(uname)" == "Darwin" ]; then
    fswatch -e ".*" -i "\\.go$" -x -r -1 .
  else
    # FILES_TO_WATCH can be something like: ./*.go ./**/*.go ./**/**/*.go
    # depending on how long the folder structure is.
    # unfortunately, inotifywait can't recursively watch specific file types
    inotifywait -e modify -e create -e delete -e move $FILES_TO_WATCH
  fi
}

build() {
  echo "Rebuilding..."
  rm -vf $TARGET
  go build -v -o $TARGET
}

kill_server() {
  if [[ -n $SERVER_PID ]]; then
    echo
    echo "Stopping server (PID: $SERVER_PID)"
    kill $SERVER_PID
  fi
}

serve() {
  echo "Starting server"
  ./$TARGET &
  SERVER_PID=$!
}

# Exit on ctrl-c (without this, ctrl-c would go to inotifywait, causing it to
# reload instead of exit):
trap "exit 0" SIGINT SIGTERM
trap kill_server "EXIT"

build
serve
while wait_for_changes; do
  kill_server
  if ! build; then
    echo "Error building server"
    continue
  fi
  serve
done
