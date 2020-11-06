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
		{"leading whitespace", " this has test-rule\n", " this has test-rule", 10, 19},
		{"no leading whitespace", "this has test-rule\n", "this has test-rule", 9, 18},
		{"leading whitespace, no new line", " this has test-rule", " this has test-rule", 10, 19},
		{"no leading whitespace, no new line", "this has test-rule", "this has test-rule", 9, 18},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			f, err := newFile(t, tc.content)
			assert.NoError(t, err)

			testRule := rule.NewTestRule()
			rules := []*rule.Rule{testRule}
			res, err := generateFileViolationsFromFilename(f.Name(), rules)
			assert.NoError(t, err)

			filename := filepath.ToSlash(f.Name())
			expected := &result.FileResults{
				Filename: filename,
				Results:  make([]result.Result, 1),
			}
			expected.Results[0] = result.LineResult{
				Rule:      rules[0],
				Violation: "test-rule",
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
		testRule := rule.NewTestRule()
		rules := []*rule.Rule{testRule}
		_, err := generateFileViolationsFromFilename("missing.file", rules)
		assert.Error(t, err)
	})

	t.Run("filename violation", func(t *testing.T) {
		f, err := newFileWithPrefix(t, "test-rule-", "content")
		assert.NoError(t, err)
		testRule := rule.NewTestRule()
		rules := []*rule.Rule{testRule}

		res, err := generateFileViolationsFromFilename(f.Name(), rules)
		assert.NoError(t, err)
		assert.Len(t, res.Results, 1)
		assert.Regexp(t, "^Filename violation: ", res.Results[0].Reason())
	})
}

// newFile creates a new file for testing. The file, and the directory that the file
// was created in will be removed at the completion of the test
func newFile(t *testing.T, text string) (*os.File, error) {
	return newFileWithPrefix(t, "woke-", text)
}

// newFile creates a new file with the prefix defined for testing.
// The file, and the directory that the file was created in will be removed at the completion of the test
func newFileWithPrefix(t *testing.T, prefix, text string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), prefix)
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
