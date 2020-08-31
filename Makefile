
build:
	docker build -t woke -f dev/Dockerfile .

run:
	docker run --rm -it \
		-v `pwd`:/go/src/github.com/get-woke/woke \
		woke "./*.go ./**/*.go ./**/**/*.go *.yaml"

.PHONY: build run

gif-record:
	go install .
	terminalizer record -k --config img/terminalizer.yml img/recording.yml
	git checkout -- test.txt
gif-play:
	terminalizer play img/recording.yml
gif-render:
	terminalizer render -o img/woke.gif img/recording.yml

.PHONY: gif-record gif-play gif-render
