
build:
	docker build -t woke -f dev/Dockerfile .

run:
	docker run --rm -it \
		-v `pwd`:/go/src/github.com/caitlinelfring/woke \
		woke "./*.go ./**/*.go ./**/**/*.go *.yaml"

.PHONY: build run
