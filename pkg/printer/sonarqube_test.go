package printer

import (
	"bytes"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestCalculateSonarSeverity(t *testing.T) {
	assert.Equal(t, "MAJOR", calculateSonarSeverity(rule.SevError))
	assert.Equal(t, "MINOR", calculateSonarSeverity(rule.SevWarn))
	assert.Equal(t, "INFO", calculateSonarSeverity(rule.SevInfo))
}

func TestSonarQube_Print(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	got := buf.String()

	expected := `{"engineId":"woke","ruleId":"whitelist","primaryLocation":{"message":"` + "`" + `whitelist` + "`" + ` may be insensitive, use ` + "`" + `allowlist` + "`" + ` instead","filePath":"foo.txt","textRange":{"startLine":1,"startColumn":6,"endColumn":15}},"type":"CODE_SMELL","severity":"MINOR"}` + "\n"
	assert.Equal(t, expected, got)
}

func TestSonarQube_PrintSuccessExitMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	assert.Equal(t, false, p.PrintSuccessExitMessage())
}

func TestSonarQube_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	p.Start()
	got := buf.String()

	expected := `{"issues":[`
	assert.Equal(t, expected, got)
}

func TestSonarQube_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	p.End()
	got := buf.String()

	expected := `]}` + "\n"
	assert.Equal(t, expected, got)
}

func TestSonarQube_Multiple(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	p.Start()
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	res = generateSecondFileResult()
	assert.NoError(t, p.Print(res))
	res = generateThirdFileResult()
	assert.NoError(t, p.Print(res))
	p.End()
	got := buf.String()

	expected := "{\"issues\":[{\"engineId\":\"woke\",\"ruleId\":\"whitelist\",\"primaryLocation\":{\"message\":\"`whitelist` may be insensitive, use `allowlist` instead\",\"filePath\":\"foo.txt\",\"textRange\":{\"startLine\":1,\"startColumn\":6,\"endColumn\":15}},\"type\":\"CODE_SMELL\",\"severity\":\"MINOR\"}\n,{\"engineId\":\"woke\",\"ruleId\":\"slave\",\"primaryLocation\":{\"message\":\"`slave` may be insensitive, use `follower` instead\",\"filePath\":\"bar.txt\",\"textRange\":{\"startLine\":1,\"startColumn\":6,\"endColumn\":15}},\"type\":\"CODE_SMELL\",\"severity\":\"MAJOR\"}\n,{\"engineId\":\"woke\",\"ruleId\":\"test\",\"primaryLocation\":{\"message\":\"`test` may be insensitive, use `alternative` instead\",\"filePath\":\"barfoo.txt\",\"textRange\":{\"startLine\":1,\"startColumn\":6,\"endColumn\":15}},\"type\":\"CODE_SMELL\",\"severity\":\"INFO\"}\n]}\n"
	assert.Equal(t, expected, got)
}
