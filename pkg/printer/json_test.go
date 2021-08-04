package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON_Print(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	got := buf.String()

	expected := "{\"Filename\":\"foo.txt\",\"Results\":[{\"Rule\":{\"Name\":\"whitelist\",\"Terms\":[\"whitelist\",\"white-list\",\"whitelisted\",\"white-listed\"],\"Alternatives\":[\"allowlist\"],\"Note\":\"\",\"Severity\":\"warning\",\"Options\":{\"WordBoundary\":false,\"WordBoundaryStart\":false,\"WordBoundaryEnd\":false,\"IncludeNote\":null}},\"Finding\":\"whitelist\",\"Line\":\"this whitelist must change\",\"StartPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":6},\"EndPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":15},\"Reason\":\"`whitelist` may be insensitive, use `allowlist` instead\"}]}\n"
	assert.Equal(t, expected, got)
}

func TestJSON_ShouldSkipExitMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	assert.Equal(t, true, p.ShouldSkipExitMessage())
}

func TestJSON_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	assert.NoError(t, p.Start())
	got := buf.String()

	expected := `{"findings": [`
	assert.Equal(t, expected, got)
}

func TestJSON_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	assert.NoError(t, p.End())
	got := buf.String()

	expected := `]}` + "\n"
	assert.Equal(t, expected, got)
}

func TestJSON_Multiple(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewJSON(buf)
	assert.NoError(t, p.Start())
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	res = generateSecondFileResult()
	assert.NoError(t, p.Print(res))
	assert.NoError(t, p.End())
	got := buf.String()

	expected := "{\"findings\": [{\"Filename\":\"foo.txt\",\"Results\":[{\"Rule\":{\"Name\":\"whitelist\",\"Terms\":[\"whitelist\",\"white-list\",\"whitelisted\",\"white-listed\"],\"Alternatives\":[\"allowlist\"],\"Note\":\"\",\"Severity\":\"warning\",\"Options\":{\"WordBoundary\":false,\"WordBoundaryStart\":false,\"WordBoundaryEnd\":false,\"IncludeNote\":null}},\"Finding\":\"whitelist\",\"Line\":\"this whitelist must change\",\"StartPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":6},\"EndPosition\":{\"Filename\":\"foo.txt\",\"Offset\":0,\"Line\":1,\"Column\":15},\"Reason\":\"`whitelist` may be insensitive, use `allowlist` instead\"}]}\n,{\"Filename\":\"bar.txt\",\"Results\":[{\"Rule\":{\"Name\":\"whitelist\",\"Terms\":[\"whitelist\",\"white-list\",\"whitelisted\",\"white-listed\"],\"Alternatives\":[\"allowlist\"],\"Note\":\"\",\"Severity\":\"warning\",\"Options\":{\"WordBoundary\":false,\"WordBoundaryStart\":false,\"WordBoundaryEnd\":false,\"IncludeNote\":null}},\"Finding\":\"blacklist\",\"Line\":\"this blacklist must change\",\"StartPosition\":{\"Filename\":\"bar.txt\",\"Offset\":0,\"Line\":1,\"Column\":6},\"EndPosition\":{\"Filename\":\"bar.txt\",\"Offset\":0,\"Line\":1,\"Column\":15},\"Reason\":\"`blacklist` may be insensitive, use `allowlist` instead\"}]}\n]}\n"
	assert.Equal(t, expected, got)
}
