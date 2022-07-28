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
	testResults := result.FileResults{
		Filename: "my/file",
		Results: []result.Result{
			result.LineResult{
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
			},
		},
	}
	buf := new(bytes.Buffer)
	p := NewCheckstyle(buf)
	p.Start()
	p.Print(&testResults)
	p.End()
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<checkstyle version="5.0">
  <file name="my/file">
    <error column="3" line="5" message="` + "`" + `whitelist` + "`" + ` may be insensitive, use ` + "`" + `allowlist` + "`" + ` instead" severity="warning" source="woke"></error>
  </file>
</checkstyle>`
	assert.Equal(t, expected, buf.String())
}

func TestCheckstyle_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewCheckstyle(buf)
	p.Start()
	got := buf.String()

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<checkstyle version="5.0">`
	assert.Equal(t, expected, got)
}

func TestCheckstyle_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewCheckstyle(buf)
	assert.PanicsWithError(t, "xml: end tag </checkstyle> without start tag", func() { p.End() })
}
