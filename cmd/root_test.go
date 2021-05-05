package cmd

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/get-woke/woke/pkg/output"
	"github.com/get-woke/woke/pkg/parser"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func TestParseArgs(t *testing.T) {
	t.Cleanup(func() {
		stdin = false
	})
	assert.Equal(t, parser.DefaultPath, parseArgs([]string{}))
	assert.Equal(t, []string{"../.."}, parseArgs([]string{"../.."}))

	stdin = true
	assert.Equal(t, []string{os.Stdin.Name()}, parseArgs([]string{}))
	assert.Equal(t, []string{os.Stdin.Name()}, parseArgs([]string{"../.."}))
}

func TestRunE(t *testing.T) {
	origStdout := output.Stdout
	origConfigFile := viper.ConfigFileUsed()
	t.Cleanup(func() {
		exitOneOnFailure = false
		noIgnore = false
		// Reset back to original
		output.Stdout = origStdout
		viper.SetConfigName(origConfigFile)
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

	t.Run("no violations found with custom message", func(t *testing.T) {
		buf := new(bytes.Buffer)
		output.Stdout = buf

		viper.SetConfigFile("../testdata/.woke-custom-exit-success.yaml")
		err := rootRunE(new(cobra.Command), []string{"../testdata/good.yml"})
		assert.NoError(t, err)

		got := buf.String()
		expected := "this is a test\n"
		assert.Equal(t, expected, got)
	})

	t.Run("violations w error", func(t *testing.T) {
		exitOneOnFailure = true

		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(`^files with violations: \d`), err.Error())
	})
}
