
FROM golang:1.17

RUN apt-get update && apt-get install -y inotify-tools

ENV ROOT_PATH /go/src/github.com/get-woke/woke
WORKDIR $ROOT_PATH
COPY go.mod ./
COPY go.sum ./

RUN go mod download

ENTRYPOINT ["./dev/autoreload.sh"]
