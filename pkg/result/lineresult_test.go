package result

import (
	"fmt"
	"go/token"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestFindResults(t *testing.T) {
	r := rule.NewTestRule()
	rs := FindResults(r, "my/file", "this has the term test-rule", 1)
	assert.Len(t, rs, 1)
	assert.Equal(t, r.Reason("test-rule"), rs[0].Reason())
	assert.Equal(t, fmt.Sprintf("    my/file:1:18-my/file:1:27 error      %s", rs[0].Reason()), rs[0].String())

	rs = FindResults(r, "my/file", "this has no rule violations", 1)
	assert.Len(t, rs, 0)

	rs = FindResults(r, "my/file", "this has the term test-rule #wokeignore:rule=test-rule", 1)
	assert.Len(t, rs, 0)
}

func TestLineResult_MarshalJSON(t *testing.T) {
	lr := testLineResult()
	b, err := lr.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(b), fmt.Sprintf(`"Reason":"%s"`, lr.Reason()))
}

func TestLineResult_GetSeverity(t *testing.T) {
	lr := testLineResult()
	assert.Equal(t, lr.GetSeverity(), lr.Rule.Severity)
}

func TestLineResult_GetStartPosition(t *testing.T) {
	lr := testLineResult()
	assert.Equal(t, lr.GetStartPosition(), lr.StartPosition)
}

func TestLineResult_GetEndPosition(t *testing.T) {
	lr := testLineResult()
	assert.Equal(t, lr.GetEndPosition(), lr.EndPosition)
}

func TestLineResult_GetLine(t *testing.T) {
	lr := testLineResult()
	assert.Equal(t, lr.GetLine(), lr.Line)
}

func testLineResult() LineResult {
	r := rule.NewTestRule()
	return LineResult{
		Rule:          r,
		Violation:     "test-rule",
		Line:          "test-rule",
		StartPosition: &token.Position{Line: 1, Offset: 0},
		EndPosition:   &token.Position{Line: 1, Offset: 8},
	}
}
