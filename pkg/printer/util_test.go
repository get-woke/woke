package printer

import (
	"bytes"
	"go/token"
	"io"
	"os"
	"sync"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
)

func generateFileResult() *result.FileResults {
	r := result.FileResults{Filename: "foo.txt"}
	r.Results = generateResults(r.Filename)
	return &r
}

func generateResults(filename string) []result.Result {
	return []result.Result{
		result.LineResult{
			Rule:      &rule.BlacklistRule,
			Violation: "blacklist",
			Line:      "this blacklist must change",
			StartPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   6,
			},
			EndPosition: &token.Position{
				Filename: filename,
				Offset:   0,
				Line:     1,
				Column:   15,
			},
		},
	}
}

// Returns output of `os.Stdout` as string.
// Based on https://medium.com/@hau12a1/golang-capturing-log-println-and-fmt-println-output-770209c791b4
func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer

	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}

func newPosition(f string, l, c int) *token.Position {
	return &token.Position{
		Filename: f,
		Offset:   0,
		Line:     l,
		Column:   c,
	}
}
