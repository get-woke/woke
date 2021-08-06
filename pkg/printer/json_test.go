package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON_Print_JSON(t *testing.T) {
	buf := new(bytes.Buffer)
	res := generateFileResult()
	p := NewJSON(buf)
	assert.NoError(t, p.Print(res))
	expected := "{\"Filename\":\"foo.txt\",\"Results\":[{\"Rule\":{\"Name\":\"whitelist\",\"Terms\":[\"whitelist\",\"white-list\",\"whitelisted\",\"white-listed\"],\"Alternatives\":[\"allowlist\"],\"Note\":\"\",\"Severity\":\"warning\",\"Options\":{\"WordBoundary\":false,\"WordBoundaryStart\":false,\"WordBoundaryEnd\":false,\"IncludeNote\":null}},\"Finding\":\"whitelist\",\"Line\":\"this whitelist must change\",\"StartPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":6},\"EndPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":15},\"Reason\":\"`whitelist` may be insensitive, use `allowlist` instead\"}]}\n"
	got := buf.String()
	assert.Equal(t, expected, got)
}

func TestJSON_ShouldSkipExitMessage(t *testing.T) {
	buf := new(bytes.Buffer)

	p := NewJSON(buf)
	assert.Equal(t, false, p.ShouldSkipExitMessage())
}

func TestJSON_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	p.Start()
	got := buf.String()

	expected := ``
	assert.Equal(t, expected, got)
}

func TestJSON_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	p.End()
	got := buf.String()

	expected := ``
	assert.Equal(t, expected, got)
}

func TestJSON_Multiple(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	p.Start()
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	res = generateSecondFileResult()
	assert.NoError(t, p.Print(res))
	res = generateThirdFileResult()
	assert.NoError(t, p.Print(res))
	p.End()
	got := buf.String()

	expected := "{\"Filename\":\"foo.txt\",\"Results\":[{\"Rule\":{\"Name\":\"whitelist\",\"Terms\":[\"whitelist\",\"white-list\",\"whitelisted\",\"white-listed\"],\"Alternatives\":[\"allowlist\"],\"Note\":\"\",\"Severity\":\"warning\",\"Options\":{\"WordBoundary\":false,\"WordBoundaryStart\":false,\"WordBoundaryEnd\":false,\"IncludeNote\":null}},\"Finding\":\"whitelist\",\"Line\":\"this whitelist must change\",\"StartPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":6},\"EndPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":15},\"Reason\":\"`whitelist` may be insensitive, use `allowlist` instead\"}]}\n{\"Filename\":\"bar.txt\",\"Results\":[{\"Rule\":{\"Name\":\"slave\",\"Terms\":[\"slave\"],\"Alternatives\":[\"follower\"],\"Note\":\"\",\"Severity\":\"error\",\"Options\":{\"WordBoundary\":false,\"WordBoundaryStart\":false,\"WordBoundaryEnd\":false,\"IncludeNote\":null}},\"Finding\":\"slave\",\"Line\":\"this slave term must change\",\"StartPosition\":{\"Filename\":\"bar.txt\",\"Offset\":0,\"Line\":1,\"Column\":6},\"EndPosition\":{\"Filename\":\"bar.txt\",\"Offset\":0,\"Line\":1,\"Column\":15},\"Reason\":\"`slave` may be insensitive, use `follower` instead\"}]}\n{\"Filename\":\"barfoo.txt\",\"Results\":[{\"Rule\":{\"Name\":\"test\",\"Terms\":[\"test\"],\"Alternatives\":[\"alternative\"],\"Note\":\"\",\"Severity\":\"info\",\"Options\":{\"WordBoundary\":false,\"WordBoundaryStart\":false,\"WordBoundaryEnd\":false,\"IncludeNote\":null}},\"Finding\":\"test\",\"Line\":\"this test must change\",\"StartPosition\":{\"Filename\":\"barfoo.txt\",\"Offset\":0,\"Line\":1,\"Column\":6},\"EndPosition\":{\"Filename\":\"barfoo.txt\",\"Offset\":0,\"Line\":1,\"Column\":15},\"Reason\":\"`test` may be insensitive, use `alternative` instead\"}]}\n"
	assert.Equal(t, expected, got)
}
