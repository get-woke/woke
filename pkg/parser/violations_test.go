package parser

import (
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFileViolations(t *testing.T) {
	tests := []struct {
		desc    string
		content string
		line    string
		start   int
		end     int
	}{
		{"leading whitespace", " this has whitelist\n", " this has whitelist", 10, 19},
		{"no leading whitespace", "this has whitelist\n", "this has whitelist", 9, 18},
		{"leading whitespace, no new line", " this has whitelist", " this has whitelist", 10, 19},
		{"no leading whitespace, no new line", "this has whitelist", "this has whitelist", 9, 18},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			f, err := newFile(t, tc.content)
			assert.NoError(t, err)

			res, err := generateFileViolationsFromFilename(f.Name(), rule.DefaultRules)
			assert.NoError(t, err)

			filename := filepath.ToSlash(f.Name())
			expected := &result.FileResults{
				Filename: filename,
				Results:  make([]result.Result, 1),
			}
			expected.Results[0] = result.LineResult{
				Rule:      &rule.WhitelistRule,
				Violation: "whitelist",
				Line:      tc.line,
				StartPosition: &token.Position{
					Filename: filename,
					Offset:   0,
					Line:     1,
					Column:   tc.start,
				},
				EndPosition: &token.Position{
					Filename: filename,
					Offset:   0,
					Line:     1,
					Column:   tc.end,
				},
			}
			assert.EqualValues(t, expected, res)
		})
	}
	t.Run("missing file", func(t *testing.T) {
		_, err := generateFileViolationsFromFilename("missing.file", rule.DefaultRules)
		assert.Error(t, err)
	})
}

// newFile creates a new file for testing. The file, and the directory that the file
// was created in will be removed at the completion of the test
func newFile(t *testing.T, text string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "woke-")
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	b := []byte(text)
	_, err = tmpFile.Write(b)

	return tmpFile, err
}
