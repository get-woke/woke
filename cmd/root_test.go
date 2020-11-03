package cmd

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/get-woke/woke/pkg/output"

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

func TestRunE(t *testing.T) {
	origStdout := output.Stdout
	t.Cleanup(func() {
		exitOneOnFailure = false
		noIgnore = false
		// Reset back to original
		output.Stdout = origStdout
	})

	t.Run("no violations found", func(t *testing.T) {
		buf := new(bytes.Buffer)
		output.Stdout = buf

		err := rootRunE(new(cobra.Command), []string{"../testdata/good.yml"})
		assert.NoError(t, err)

		got := buf.String()
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
