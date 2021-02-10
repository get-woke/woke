package printer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/get-woke/woke/pkg/result"

	"github.com/stretchr/testify/assert"
)

func TestText_Print(t *testing.T) {
	p := NewText(true)
	res := generateFileResult()
	buf := new(bytes.Buffer)
	assert.NoError(t, p.Print(buf, res))
	got := buf.String()
	expected := fmt.Sprintf("foo.txt:1:6-15: %s (%s)\n%s\n      ^\n", res.Results[0].Reason(), res.Results[0].GetSeverity(), res.Results[0].GetLine())
	assert.Equal(t, expected, got)
}

func TestText_arrowUnderLine(t *testing.T) {
	p := NewText(true)

	r := result.LineResult{
		Line:          "this line has test-rule as a violation",
		StartPosition: newPosition("foo.txt", 4, 14),
		EndPosition:   newPosition("foo.txt", 4, 24),
	}
	assert.Equal(t, "              ^", p.arrowUnderLine(&r))

	r = result.LineResult{
		Line:          "    this line has test-rule as a violation",
		StartPosition: newPosition("foo.txt", 4, 18),
		EndPosition:   newPosition("foo.txt", 4, 28),
	}
	assert.Equal(t, "                  ^", p.arrowUnderLine(&r))

	r = result.LineResult{
		Line:          "\tthis line has test-rule as a violation",
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
