package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/get-woke/woke/pkg/output"
	"github.com/get-woke/woke/pkg/parser"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// run profiling with
// go test -v -cpuprofile cpu.prof -memprofile mem.prof -bench=. ./cmd
// memory:
//
//	go tool pprof mem.prof
//
// cpu:
//
//	go tool pprof cpu.prof
func BenchmarkRootRunE(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
	output.Stdout = io.Discard
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assert.NoError(b, rootRunE(new(cobra.Command), []string{".."}))
	}
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

	overrideHomeDir(t)

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
	tests := []struct {
		stdin         bool
		args          []string
		expectedArgs  []string
		expectedError error
	}{
		{
			stdin:         false,
			args:          []string{},
			expectedArgs:  parser.DefaultPath,
			expectedError: nil,
		},
		{
			stdin:         false,
			args:          []string{"../.."},
			expectedArgs:  []string{filepath.Join("..", "..")},
			expectedError: nil,
		},

		// Test glob expansion
		{
			stdin: false,
			args:  []string{"../testdata/*.yml"},
			expectedArgs: []string{
				filepath.Join("..", "testdata/bad.yml"),
				filepath.Join("..", "testdata/good.yml"),
			},
			expectedError: nil,
		},
		{
			stdin:         false,
			args:          []string{"../testdata/g??d.yml"}, // matches any single non-separator character
			expectedArgs:  []string{filepath.Join("..", "testdata/good.yml")},
			expectedError: nil,
		},
		{
			stdin:         false,
			args:          []string{"../testdata/[a-z]ood.yml"}, // character range
			expectedArgs:  []string{filepath.Join("..", "testdata", "good.yml")},
			expectedError: nil,
		},
		{
			stdin:         false,
			args:          []string{"../testdata/[^abc]ood.yml"}, // character class with negation.
			expectedArgs:  []string{filepath.Join("..", "testdata", "good.yml")},
			expectedError: nil,
		},
		{
			stdin:         false,
			args:          []string{"../testdata/[!abc]ood.yml"}, // character class with negation.
			expectedArgs:  []string{filepath.Join("..", "testdata", "good.yml")},
			expectedError: nil,
		},
		{
			stdin:         false,
			args:          []string{"../testdata/[^g]ood.yml"}, // character class with negation.
			expectedArgs:  nil,
			expectedError: nil,
		},
		{
			stdin: false,
			args:  []string{"../testdata/*/*.yml"},
			expectedArgs: []string{
				filepath.Join("..", "testdata", "subdir1", "bad.yml"),
				filepath.Join("..", "testdata", "subdir1", "good.yml"),
			},
			expectedError: nil,
		},
		{
			stdin: false,
			args:  []string{"../testdata/**/*.yml"},
			expectedArgs: []string{
				filepath.Join("..", "testdata", "bad.yml"),
				filepath.Join("..", "testdata", "good.yml"),
				filepath.Join("..", "testdata", "subdir1", "bad.yml"),
				filepath.Join("..", "testdata", "subdir1", "good.yml"),
				filepath.Join("..", "testdata", "subdir1", "subdir2", "bad.yml"),
				filepath.Join("..", "testdata", "subdir1", "subdir2", "good.yml"),
			},
			expectedError: nil,
		},
		{
			stdin: false,
			args:  []string{"../testdata/**/{good,bad}.yml"}, // Alternate pattern
			expectedArgs: []string{
				filepath.Join("..", "testdata", "bad.yml"),
				filepath.Join("..", "testdata", "good.yml"),
				filepath.Join("..", "testdata", "subdir1", "bad.yml"),
				filepath.Join("..", "testdata", "subdir1", "good.yml"),
				filepath.Join("..", "testdata", "subdir1", "subdir2", "bad.yml"),
				filepath.Join("..", "testdata", "subdir1", "subdir2", "good.yml"),
			},
			expectedError: nil,
		},
		{
			stdin: false,
			args:  []string{"../testdata/**/?ood.yml"},
			expectedArgs: []string{
				filepath.Join("..", "testdata", "good.yml"),
				filepath.Join("..", "testdata", "subdir1", "good.yml"),
				filepath.Join("..", "testdata", "subdir1", "subdir2", "good.yml"),
			},
			expectedError: nil,
		},

		// Bad glob pattern
		{
			stdin:         false,
			args:          []string{"r[.go"}, // Invalid character class
			expectedArgs:  nil,
			expectedError: doublestar.ErrBadPattern,
		},
		{
			stdin:         false,
			args:          []string{"{.go"}, // Bad alternate pattern
			expectedArgs:  nil,
			expectedError: doublestar.ErrBadPattern,
		},

		{
			stdin:         true,
			args:          []string{},
			expectedArgs:  []string{os.Stdin.Name()},
			expectedError: nil,
		},
		{
			stdin:         true,
			args:          []string{"../.."},
			expectedArgs:  []string{os.Stdin.Name()},
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			stdin = tt.stdin
			files, err := parseArgs(tt.args)
			assert.ErrorIs(t, err, tt.expectedError,
				fmt.Sprintf("arguments: %v. Expected '%v', Got '%v'", tt.args, err, tt.expectedError))
			assert.Equal(t, tt.expectedArgs, files)
		})
	}
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

	t.Run("findings with inclusive language issues", func(t *testing.T) {
		exitOneOnFailure = true
		// don't ignore testdata folder
		noIgnore = true

		t.Cleanup(func() {
			exitOneOnFailure = false
		})
		err := rootRunE(new(cobra.Command), []string{"../testdata"})
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(`^files with findings: \d`), err.Error())
	})

	t.Run("findings with invalid glob pattern", func(t *testing.T) {
		exitOneOnFailure = true
		// don't ignore testdata folder
		noIgnore = true

		t.Cleanup(func() {
			exitOneOnFailure = false
		})
		err := rootRunE(new(cobra.Command), []string{"../testdata/**/[.yml"})
		assert.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(`syntax error in pattern`), err.Error())
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

// helper functions

func setTestConfigFile(t *testing.T, filename string) {
	origConfigFile := viper.ConfigFileUsed()
	t.Cleanup(func() {
		viper.SetConfigFile(origConfigFile)
	})
	viper.SetConfigFile(filename)
}

// overrideHomeDir to avoid pulling in a config file in the home directory
// while running tests
func overrideHomeDir(t *testing.T) {
	origHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", origHome)
		homedir.Reset()
	})
	os.Setenv("HOME", "foo")
	homedir.Reset()
}
