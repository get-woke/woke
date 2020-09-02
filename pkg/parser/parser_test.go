package parser

import (
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestParser_ParsePaths(t *testing.T) {
	i, err := ignore.NewIgnore()
	assert.NoError(t, err)
	p := NewParser(rule.DefaultRules, i)

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

	// Test ignored file
	f4, err := newFile("i have a whitelist violation, but am ignored\n")
	assert.NoError(t, err)
	defer os.Remove(f4.Name())
	i2, err := ignore.NewIgnore(f4.Name())
	assert.NoError(t, err)
	p.Ignorer = i2
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
