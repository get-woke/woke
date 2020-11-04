package printer

import (
	"fmt"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestFormatResultForGitHubAction(t *testing.T) {
	res := generateFileResult().Results[0]

	got := formatResultForGitHubAction(res)
	assert.Equal(t, "::error file=foo.txt,line=1,col=6::"+res.Reason(), got)
}

func TestTranslateSeverityForAction(t *testing.T) {
	assert.Equal(t, translateSeverityForAction(rule.SevError), "error")
	assert.Equal(t, translateSeverityForAction(rule.SevWarn), "warning")
	assert.Equal(t, translateSeverityForAction(rule.SevInfo), "warning")
}

func TestGitHubActions_Print(t *testing.T) {
	p := NewGitHubActions()
	res := generateFileResult()
	got := captureOutput(func() {
		assert.NoError(t, p.Print(res))
	})
	expected := fmt.Sprintf("::error file=foo.txt,line=1,col=6::%s\n", res.Results[0].Reason())
	assert.Equal(t, expected, got)
}
