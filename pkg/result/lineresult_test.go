package result

import (
	"fmt"
	"go/token"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestFindResults(t *testing.T) {
	rs := FindResults(&rule.TestRule, "my/file", "this has the term whitelist", 1)
	assert.Len(t, rs, 1)
	assert.Equal(t, rule.TestRule.Reason("whitelist"), rs[0].Reason())
	assert.Equal(t, fmt.Sprintf("    my/file:1:18-my/file:1:27 warning    %s", rs[0].Reason()), rs[0].String())

	rs = FindResults(&rule.TestRule, "my/file", "this has no rule findings", 1)
	assert.Len(t, rs, 0)

	// inline-ignoring is handled in Parser.generateFileFindings, not FindResults
	rs = FindResults(&rule.TestRule, "my/file", "this has the term whitelist #wokeignore:rule=whitelist", 1)
	assert.Len(t, rs, 1)
	rs = FindResults(&rule.TestRule, "my/file", "/* wokeignore:rule=whitelist */ this has the term whitelist", 1)
	assert.Len(t, rs, 1)
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

func TestLineResult_GetRuleName(t *testing.T) {
	lr := testLineResult()
	assert.Equal(t, lr.GetRuleName(), lr.Rule.Name)
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
	return LineResult{
		Rule:          &rule.TestRule,
		Finding:       "whitelist",
		Line:          "whitelist",
		StartPosition: &token.Position{Line: 1, Offset: 0},
		EndPosition:   &token.Position{Line: 1, Offset: 8},
	}
}
