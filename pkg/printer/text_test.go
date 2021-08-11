package printer

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/get-woke/woke/pkg/result"

	"github.com/stretchr/testify/assert"
)

func TestText_Print(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewText(buf, true)
	res := generateFileResult()
	assert.NoError(t, p.Print(res))
	got := buf.String()
	expected := fmt.Sprintf("foo.txt:1:6-15: %s (%s)\n%s\n      ^\n", res.Results[0].Reason(), res.Results[0].GetSeverity(), res.Results[0].GetLine())
	assert.Equal(t, expected, got)
}

func TestText_Start(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewText(buf, true)
	p.Start()
	got := buf.String()
	assert.Equal(t, ``, got)
}

func TestText_End(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewText(buf, true)
	p.End()
	got := buf.String()
	assert.Equal(t, ``, got)
}

func TestText_PrintSuccessExitMessage(t *testing.T) {
	buf := new(bytes.Buffer)
	p := NewText(buf, true)
	assert.Equal(t, true, p.PrintSuccessExitMessage())
}

func TestText_arrowUnderLine(t *testing.T) {
	p := NewText(io.Discard, true)

	r := result.LineResult{
		Line:          "this line has black-list as a finding",
		StartPosition: newPosition("foo.txt", 4, 14),
		EndPosition:   newPosition("foo.txt", 4, 24),
	}
	assert.Equal(t, "              ^", p.arrowUnderLine(&r))

	r = result.LineResult{
		Line:          "    this line has black-list as a finding",
		StartPosition: newPosition("foo.txt", 4, 18),
		EndPosition:   newPosition("foo.txt", 4, 28),
	}
	assert.Equal(t, "                  ^", p.arrowUnderLine(&r))

	r = result.LineResult{
		Line:          "\tthis line has black-list as a finding",
		StartPosition: newPosition("foo.txt", 4, 15),
		EndPosition:   newPosition("foo.txt", 4, 25),
	}
	assert.Equal(t, "\t              ^", p.arrowUnderLine(&r))

	r = result.LineResult{
		Line:          "unknown",
		StartPosition: newPosition("foo.txt", 1, 0),
		EndPosition:   newPosition("foo.txt", 1, 0),
	}
	assert.Equal(t, "", p.arrowUnderLine(&r))
}
