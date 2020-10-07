package cmd

import (
	"regexp"
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

func TestRunE(t *testing.T) {
	t.Cleanup(func() {
		exitOneOnFailure = false
		noIgnore = false
	})
	t.Run("no violations", func(t *testing.T) {
		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.NoError(t, err)
	})
	t.Run("violations w error", func(t *testing.T) {
		exitOneOnFailure = true
		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(`^files with violations: \d`), err.Error())
	})
}
