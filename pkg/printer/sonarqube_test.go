package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSonarQube_Print(t *testing.T) {
	p := NewSonarQube()
	res := generateFileResult()
	buf := new(bytes.Buffer)
	assert.NoError(t, p.Print(buf, res))
	got := buf.String()

	expected := `[{"engineId":"woke","ruleId":"whitelist","primaryLocation":{"message":"` + "`" + `whitelist` + "`" + ` may be insensitive, use ` + "`" + `allowlist` + "`" + ` instead","filePath":"foo.txt","textRange":{"startLine":1,"startColumn":6,"endColumn":15}},"type":"CODE_SMELL","severity":"MAJOR"}]` + "\n"
	assert.Equal(t, expected, got)
}
