
build:
	docker build -t woke -f dev/Dockerfile .

run:
	docker run --rm -it \
		-v `pwd`:/go/src/github.com/get-woke/woke \
		woke "./*.go ./**/*.go ./**/**/*.go *.yaml"

.PHONY: build run

prof:
	go test -bench=. -run=^$$ -cpuprofile cpu.prof -memprofile mem.prof ./cmd

prof-mem: prof
	pprof -top mem.prof | head -n 10

# pprof -http=localhost:8080 mem.prof

prof-cpu: prof
	pprof -http=localhost:8080 cpu.prof

.PHONY: prof prof-mem prof-cpu
