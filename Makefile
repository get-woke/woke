
build:
	docker build -t woke -f Dockerfile.dev .

run:
	docker run --rm -it \
		-v `pwd`:/go/src/github.com/caitlinelfring/woke \
		woke bin/woke "./*.go ./**/*.go ./**/**/*.go"

.PHONY: build run
