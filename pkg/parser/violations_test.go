package parser

import (
	"go/token"
	"strings"
	"testing"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFileViolations(t *testing.T) {
	violatingPrefix := "whitelist"
	f, err := newFileWithPrefix(t, violatingPrefix+"-", "this has whitelist\n")
	assert.NoError(t, err)
	column := strings.Index(f.Name(), violatingPrefix)

	res, err := generateFileViolationsFromFilename(f.Name(), rule.DefaultRules)
	assert.NoError(t, err)

	expected := &result.FileResults{
		Filename: f.Name(),
		Results: []result.Result{
			{
				Rule:      &rule.WhitelistRule,
				Violation: "whitelist",
				Line:      f.Name(),
				StartPosition: &token.Position{
					Filename: f.Name(),
					Offset:   0,
					Line:     1,
					Column:   column,
				},
				EndPosition: &token.Position{
					Filename: f.Name(),
					Offset:   0,
					Line:     1,
					Column:   column + len(violatingPrefix),
				},
			},
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
