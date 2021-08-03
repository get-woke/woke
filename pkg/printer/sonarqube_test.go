package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSonarQube_Print(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	got := buf.String()

	expected := `{"engineId":"woke","ruleId":"whitelist","primaryLocation":{"message":"` + "`" + `whitelist` + "`" + ` may be insensitive, use ` + "`" + `allowlist` + "`" + ` instead","filePath":"foo.txt","textRange":{"startLine":1,"startColumn":6,"endColumn":15}},"type":"CODE_SMELL","severity":"MAJOR"}` + "\n"
	assert.Equal(t, expected, got)
}

func TestSonarQube_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	assert.NoError(t, p.Start())
	got := buf.String()

	expected := `{"issues":[`
	assert.Equal(t, expected, got)
}

func TestSonarQube_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewSonarQube(buf)
	assert.NoError(t, p.End())
	got := buf.String()

	expected := `]}` + "\n"
	assert.Equal(t, expected, got)
}
