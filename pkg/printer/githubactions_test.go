package printer

import (
	"go/token"
	"testing"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestFormatResultForGitHubAction(t *testing.T) {
	testResult := result.Result{
		Rule:      &rule.WhitelistRule,
		Violation: "whitelist",
		StartPosition: &token.Position{
			Filename: "my/file",
			Offset:   0,
			Line:     5,
			Column:   3,
		},
		EndPosition: &token.Position{
			Filename: "my/file",
			Offset:   0,
			Line:     5,
			Column:   12,
		},
	}
	got := formatResultForGitHubAction(&testResult)
	assert.Equal(t, "::warning file=my/file,line=5,col=3::"+testResult.Rule.Reason(testResult.Violation), got)
}

func TestTranslateSeverityForAction(t *testing.T) {
	assert.Equal(t, translateSeverityForAction(rule.SevError), "error")
	assert.Equal(t, translateSeverityForAction(rule.SevWarn), "warning")
	assert.Equal(t, translateSeverityForAction(rule.SevInfo), "warning")
}

func TestGitHubActions_Print(t *testing.T) {
	p := NewGitHubActions()
	got := captureOutput(func() {
		assert.NoError(t, p.Print(generateFileResult()))
	})
	expected := "::warning file=foo.txt,line=5,col=3::`blacklist` may be insensitive, use `blocklist` instead\n"
	assert.Equal(t, expected, got)
}
