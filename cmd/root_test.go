package cmd

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/get-woke/woke/pkg/output"
	"github.com/get-woke/woke/pkg/parser"

	"github.com/rs/zerolog"
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
func BenchmarkRootRunE(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
	output.Stdout = io.Discard
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assert.NoError(b, rootRunE(new(cobra.Command), []string{".."}))
	}
}

func setTestConfigFile(t *testing.T, filename string) {
	origConfigFile := viper.ConfigFileUsed()
	t.Cleanup(func() {
		viper.SetConfigFile(origConfigFile)
	})
	viper.SetConfigFile(filename)
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		desc    string
		cfgFile string
	}{
		{
			desc:    "good config",
			cfgFile: "../testdata/good.yml",
		},
		{
			desc:    "no config",
			cfgFile: "",
		},
		{
			desc:    "invalid config",
			cfgFile: "../testdata/invalid.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Cleanup(func() {
				cfgFile = ""
				initConfig()
			})
			cfgFile = tt.cfgFile
			initConfig()
			assert.Equal(t, tt.cfgFile, viper.ConfigFileUsed())
		})
	}
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
	t.Cleanup(func() {
		// Reset back to original
		output.Stdout = origStdout
	})

	t.Run("no findings found", func(t *testing.T) {
		buf := new(bytes.Buffer)
		output.Stdout = buf

		err := rootRunE(new(cobra.Command), []string{"../testdata/good.yml"})
		assert.NoError(t, err)

		got := buf.String()
		expected := "No findings found.\n"
		assert.Equal(t, expected, got)
	})

	t.Run("no findings found with custom message", func(t *testing.T) {
		buf := new(bytes.Buffer)
		output.Stdout = buf
		setTestConfigFile(t, "../testdata/.woke-custom-exit-success.yaml")
		err := rootRunE(new(cobra.Command), []string{"../testdata/good.yml"})
		assert.NoError(t, err)

		got := buf.String()
		expected := "this is a test\n"
		assert.Equal(t, expected, got)
	})

	t.Run("findings w error", func(t *testing.T) {
		exitOneOnFailure = true
		t.Cleanup(func() {
			exitOneOnFailure = false
		})
		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(`^files with findings: \d`), err.Error())
	})

	t.Run("no rules enabled", func(t *testing.T) {
		disableDefaultRules = true
		t.Cleanup(func() {
			disableDefaultRules = false
		})

		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNoRulesEnabled)
	})

	t.Run("invalid printer", func(t *testing.T) {
		outputName = "foo"
		t.Cleanup(func() {
			outputName = "text"
		})
		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
		assert.Equal(t, "foo is not a valid printer type", err.Error())
	})

	t.Run("invalid config", func(t *testing.T) {
		setTestConfigFile(t, "../testdata/invalid.yaml")
		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
	})
}
