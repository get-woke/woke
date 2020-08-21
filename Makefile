
build:
	docker build -t woke -f Dockerfile.dev .

run:
	docker run --rm -it \
		-v `pwd`:/go/src/github.com/caitlinelfring/woke \
		woke "./*.go ./**/*.go ./**/**/*.go *.yaml"

.PHONY: build run
