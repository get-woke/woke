package parser

import (
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func parsePathTests(t *testing.T) {
	p := NewParser(rule.DefaultRules, ignore.NewIgnore([]string{}, []string{}))

	f1, err := newFile("i have a whitelist\n")
	assert.NoError(t, err)
	defer os.Remove(f1.Name())

	fr1, err := p.ParsePaths(f1.Name())
	assert.NoError(t, err)

	assert.Len(t, fr1, 1)
	assert.EqualValues(t, result.FileResults{
		Filename: f1.Name(),
		Results: []result.Result{
			{
				Rule:      &rule.WhitelistRule,
				Violation: "whitelist",
				Line:      "i have a whitelist",
				StartPosition: &token.Position{
					Filename: f1.Name(),
					Offset:   0,
					Line:     1,
					Column:   9,
				},
				EndPosition: &token.Position{
					Filename: f1.Name(),
					Offset:   0,
					Line:     1,
					Column:   18,
				},
			},
		},
	}, fr1[0])

	f2, err := newFile("i have a no violations\n")
	assert.NoError(t, err)
	defer os.Remove(f2.Name())
	fr2, err := p.ParsePaths(f2.Name())
	assert.NoError(t, err)
	assert.Len(t, fr2, 0)

	// Test for IsTextFileFromFilename failure
	f3, err := newFile("")
	assert.NoError(t, err)
	defer os.Remove(f3.Name())
	fr3, err := p.ParsePaths(f3.Name())
	assert.NoError(t, err)
	assert.Len(t, fr3, 0)

	// Test with multiple paths supplied
	// Disabled since this functionality is currently broken
	frAll, err := p.ParsePaths(f1.Name(), f2.Name(), f3.Name())
	assert.NoError(t, err)
	assert.Len(t, frAll, 1)

	// Test ignored file
	f4, err := newFile("i have a whitelist violation, but am ignored\n")
	assert.NoError(t, err)
	defer os.Remove(f4.Name())

	p.Ignorer = ignore.NewIgnore([]string{f4.Name()}, []string{})
	fr4, err := p.ParsePaths(f4.Name())
	assert.NoError(t, err)
	assert.Len(t, fr4, 0)

	// Test default path (which would run tests against the parser package)
	fr5, err := p.ParsePaths()
	assert.NoError(t, err)
	assert.Greater(t, len(fr5), 0)

	// Stdin
	err = writeToStdin("i have a whitelist here\n", func() {
		fr, err := p.ParsePaths(os.Stdin.Name())
		assert.NoError(t, err)
		assert.Len(t, fr, 1)
		assert.EqualValues(t, result.FileResults{
			Filename: os.Stdin.Name(),
			Results: []result.Result{
				{
					Rule:      &rule.WhitelistRule,
					Violation: "whitelist",
					Line:      "i have a whitelist here",
					StartPosition: &token.Position{
						Filename: os.Stdin.Name(),
						Offset:   0,
						Line:     1,
						Column:   9,
					},
					EndPosition: &token.Position{
						Filename: os.Stdin.Name(),
						Offset:   0,
						Line:     1,
						Column:   18,
					},
				},
			},
		}, fr[0])
	})
	assert.NoError(t, err)
}

func TestParser_ParsePaths(t *testing.T) {
	os.Unsetenv("WORKER_POOL_COUNT")
	parsePathTests(t)

	os.Setenv("WORKER_POOL_COUNT", "10")
	defer os.Unsetenv("WORKER_POOL_COUNT")
	parsePathTests(t)
}

func TestPathsIncludeStdin(t *testing.T) {
	testCases := []struct {
		paths    []string
		expected bool
	}{
		{
			paths:    []string{"file1", "file2"},
			expected: false,
		},
		{
			paths:    []string{},
			expected: false,
		},
		{
			paths:    []string{"/dev/stdin", "file2"},
			expected: true,
		},
		{
			paths:    []string{"/dev/stdin"},
			expected: true,
		},
	}
	for _, tC := range testCases {
		assert.Equal(t, tC.expected, pathsIncludeStdin(tC.paths), tC.paths)
	}
}

func writeToStdin(text string, f func()) error {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(text)); err != nil {
		return err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return err
	}

	stdin := os.Stdin
	defer func() {
		os.Stdin = stdin
	}()
	os.Stdin = tmpfile
	f()
	return tmpfile.Close()
}

func BenchmarkParsePaths(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
	tmpFile, err := ioutil.TempFile(os.TempDir(), "")
	assert.NoError(b, err)

	// Remember to clean up the file afterwards
	defer os.Remove(tmpFile.Name())

	for i := 0; i < 100; i++ {
		_, _ = tmpFile.WriteString("this whitelist, he put in man hours to sanity-check the master/slave dummy-value. we can do better.\n")
	}
	tmpFile.Close()

	for i := 0; i < b.N; i++ {
		ignorer := ignore.NewIgnore([]string{}, []string{})
		p := NewParser(rule.DefaultRules, ignorer)
		_, err = p.ParsePaths(tmpFile.Name())
		assert.NoError(b, err)
	}
}

func BenchmarkParsePathsRoot(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)

	for i := 0; i < b.N; i++ {
		Ignorer := ignore.NewIgnore([]string{}, []string{})
		p := NewParser(rule.DefaultRules, Ignorer)
		_, err := p.ParsePaths("../..")
		assert.NoError(b, err)
	}
}
