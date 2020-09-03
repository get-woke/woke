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
	f, err := newFile("this has whitelist\n")
	defer os.Remove(f.Name())

	assert.NoError(t, err)

	res, err := generateFileViolationsFromFilename(f.Name(), rule.DefaultRules)
	assert.NoError(t, err)

	expected := &result.FileResults{
		Filename: f.Name(),
		Results: []result.Result{
			{
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
			},
		},
	}
	assert.EqualValues(t, expected, res)

	_, err = generateFileViolationsFromFilename("missing.file", rule.DefaultRules)
	assert.Error(t, err)
}

func newFile(text string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "woke-")
	if err != nil {
		return nil, err
	}

	b := []byte(text)
	_, err = tmpFile.Write(b)

	return tmpFile, err

}
