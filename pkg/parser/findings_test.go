package parser

import (
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFileFindings(t *testing.T) {
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
			p, err := testParser()
			assert.NoError(t, err)
			res, err := p.generateFileFindingsFromFilename(f.Name())
			assert.NoError(t, err)

			filename := filepath.ToSlash(f.Name())
			expected := &result.FileResults{
				Filename: filename,
				Results:  make([]result.Result, 1),
			}
			expected.Results[0] = result.LineResult{
				Rule:    &rule.TestRule,
				Finding: "whitelist",
				Line:    tc.line,
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
		p, err := testParser()
		assert.NoError(t, err)
		_, err = p.generateFileFindingsFromFilename("missing.file")
		assert.Error(t, err)
	})

	t.Run("filename finding", func(t *testing.T) {
		f, err := newFileWithPrefix(t, "whitelist-", "content")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		res, err := p.generateFileFindingsFromFilename(f.Name())
		assert.NoError(t, err)
		assert.Len(t, res.Results, 1)
		assert.Regexp(t, "^Filename finding: ", res.Results[0].Reason())
	})

	t.Run("filename finding for empty file", func(t *testing.T) {
		f, err := newFileWithPrefix(t, "empty-whitelist-", "")
		assert.NoError(t, err)

		p, err := testParser()
		assert.NoError(t, err)
		res, err := p.generateFileFindingsFromFilename(f.Name())
		assert.NoError(t, err)
		assert.Len(t, res.Results, 1)
		assert.Regexp(t, "^Filename finding: ", res.Results[0].Reason())
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
	tmpFile, err := os.CreateTemp(os.TempDir(), prefix)
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

// Tests for when a rule name is used with inline wokeignore, when that rule name matches another rule
func TestGenerateFileFindingsOverlappingRules(t *testing.T) {
	tests := []struct {
		desc    string
		content string
		matches int
	}{
		{"overlapping rule", "this has master #wokeignore:rule=master-slave", 0},
		{"overlapping rule two ignores", "this has master #wokeignore:rule=master-slave,slave", 0},
		{"wrong rule", "this has whitelist # wokeignore:rule=blacklist", 1},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			f, err := newFile(t, tc.content)
			assert.NoError(t, err)

			p, err := testParser()
			assert.NoError(t, err)
			res, err := p.generateFileFindingsFromFilename(f.Name())
			assert.NoError(t, err)
			assert.Len(t, res.Results, tc.matches)
		})
	}
}

// Tests for next-line wokeignore
func TestGenerateFileFindingsNewLineIgnores(t *testing.T) {
	tests := []struct {
		desc    string
		content string
		matches int
	}{
		{"not matching newline ignore", "#wokeignore:rule=master-slave\n this has whitelist", 1},
		{"matching newline ignore", "#wokeignore:rule=whitelist\n this has whitelist", 0},
		{"matching newline ignore", "#wokeignore:rule=whitelist whitelist\n this has whitelist", 0},
		{"newline ignore with potential match two lines down", "#wokeignore:rule=whitelist\n this line is fine\n this has whitelist", 1},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			f, err := newFile(t, tc.content)
			assert.NoError(t, err)

			p, err := testParser()
			assert.NoError(t, err)
			res, err := p.generateFileFindingsFromFilename(f.Name())
			assert.NoError(t, err)
			assert.Len(t, res.Results, tc.matches)
		})
	}
}
