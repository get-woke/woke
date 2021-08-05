package printer

import (
	"bytes"
	"fmt"
	"go/token"
	"testing"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestFormatResultForGitHubAction(t *testing.T) {
	testResult := result.LineResult{
		Rule:    &rule.TestRule,
		Finding: "whitelist",
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
	assert.Equal(t, "::warning file=my/file,line=5,col=3::"+testResult.Rule.Reason(testResult.Finding), got)
}

func TestTranslateSeverityForAction(t *testing.T) {
	assert.Equal(t, translateSeverityForAction(rule.SevError), "error")
	assert.Equal(t, translateSeverityForAction(rule.SevWarn), "warning")
	assert.Equal(t, translateSeverityForAction(rule.SevInfo), "warning")
}

func TestGitHubActions_ShouldSkipExitMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewGitHubActions(buf)
	assert.Equal(t, false, p.ShouldSkipExitMessage())
}

func TestGitHubActions_Print(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewGitHubActions(buf)
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	got := buf.String()
	expected := fmt.Sprintf("::warning file=foo.txt,line=1,col=6::%s\n", res.Results[0].Reason())
	assert.Equal(t, expected, got)
}

func TestGitHubActions_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewGitHubActions(buf)
	p.Start()
	got := buf.String()
	assert.Equal(t, ``, got)
}

func TestGitHubActions_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewGitHubActions(buf)
	p.End()
	got := buf.String()
	assert.Equal(t, ``, got)
}
