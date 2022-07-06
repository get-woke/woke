package printer

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestFormatResultForCheckstyle(t *testing.T) {
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
	got := formatResultForCheckstyle(&testResult)
	assert.Equal(t,
		"    <error column=\"3\" line=\"5\" message=\"`whitelist` may be insensitive, use `allowlist` instead\" severity=\"warning\" source=\"woke\"/>",
		got)
}

func TestCheckstyle_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewCheckstyle(buf)
	p.Start()
	got := buf.String()

	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<checkstyle version=\"5.0\">\n"
	assert.Equal(t, expected, got)
}

func TestCheckstyle_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewCheckstyle(buf)
	p.End()
	got := buf.String()

	expected := "</checkstyle>\n"
	assert.Equal(t, expected, got)
}
