package parser

import (
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFileViolations(t *testing.T) {
	f, err := newFile(t, "this has whitelist\n")
	assert.NoError(t, err)

	res, err := generateFileViolationsFromFilename(f.Name(), rule.DefaultRules)
	assert.NoError(t, err)

	expected := &result.FileResults{
		Filename: f.Name(),
		Results:  make([]result.Result, 1),
	}
	expected.Results[0] = result.LineResult{
		Rule:      &rule.WhitelistRule,
		Violation: "whitelist",
		Line:      "this has whitelist",
		StartPosition: &token.Position{
			Filename: f.Name(),
			Offset:   0,
			Line:     1,
			Column:   9,
		},
		EndPosition: &token.Position{
			Filename: f.Name(),
			Offset:   0,
			Line:     1,
			Column:   18,
		},
	}
	assert.EqualValues(t, expected, res)

	_, err = generateFileViolationsFromFilename("missing.file", rule.DefaultRules)
	assert.Error(t, err)
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
