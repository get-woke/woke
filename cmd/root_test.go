package cmd

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"sync"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// run profiling with
// go test -v -cpuprofile cpu.prof -memprofile mem.prof -bench=. ./cmd
// memory:
//    go tool pprof mem.prof
// cpu:
//    go tool pprof cpu.prof
func BenchmarkExecute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		assert.NoError(b, rootRunE(new(cobra.Command), []string{".."}))
	}
}

func TestInitConfig(t *testing.T) {
	t.Cleanup(func() {
		cfgFile = ""
		debug = false
	})
	debug = true
	t.Run("good config", func(t *testing.T) {
		cfgFile = "../testdata/good.yml"
		initConfig()
	})

	t.Run("no config", func(t *testing.T) {
		cfgFile = ""
		initConfig()
	})
}

func TestRunE(t *testing.T) {
	t.Cleanup(func() {
		exitOneOnFailure = false
		noIgnore = false
	})
	t.Run("no violations found", func(t *testing.T) {
		got := captureOutput(func() {
			err := rootRunE(new(cobra.Command), []string{"../testdata/good.yml"})
			assert.NoError(t, err)
		})
		expected := "No violations found. Stay woke \u270a\n"
		assert.Equal(t, expected, got)
	})
	t.Run("violations w error", func(t *testing.T) {
		exitOneOnFailure = true
		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(`^files with violations: \d`), err.Error())
	})
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
