FROM golang:1.14 as builder

ARG BUILD_TIME=0
ARG BUILD_VERSION=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Make it runnable on alpine
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X 'main.BuildTime=${BUILD_TIME}' -X 'main.BuildVersion=${BUILD_VERSION}'" \
  -a -installsuffix cgo -o woke .

######################

FROM alpine:latest
COPY --from=builder /app/woke /woke
COPY default.yaml /default.yaml
ENTRYPOINT [ "/woke" ]
