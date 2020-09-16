package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRule_FindMatchIndexes(t *testing.T) {
	r := testRule()
	tests := []struct {
		text     string
		expected [][]int
	}{
		{"this string has rule-1 and rule1 included", [][]int{{16, 22}, {27, 32}}},
		{"this string has rule-2 and rule1 included", [][]int{{27, 32}}},
		{"this string does not have any violations", [][]int(nil)},
		{"this string has no violation due to word boundary rule1rule-1", [][]int(nil)},
		// ^ This case should be a violation, but the code needs to be updated to support it.
	}
	for _, test := range tests {
		got := r.FindMatchIndexes(test.text)
		assert.Equal(t, test.expected, got)
	}

	e := Rule{Name: "rule1"}
	assert.Equal(t, [][]int(nil), e.FindMatchIndexes("rule1"))
}

func TestRule_Reason(t *testing.T) {
	r := testRule()
	assert.Equal(t, "`rule-1` may be insensitive, use `alt-rule1`, `alt-rule-1` instead", r.Reason("rule-1"))

	r.Alternatives = []string(nil)
	assert.Equal(t, "`rule-1` may be insensitive, try not to use it", r.Reason("rule-1"))
}

func TestRule_ReasonWithNode(t *testing.T) {
	r := testRule()

	assert.Equal(t, "`rule-1` may be insensitive, use `alt-rule1`, `alt-rule-1` instead", r.ReasonWithNote("rule-1"))

	r.Note = "rule note here"
	assert.Equal(t, "`rule-1` may be insensitive, use `alt-rule1`, `alt-rule-1` instead (rule note here)", r.ReasonWithNote("rule-1"))
}

func testRule() Rule {
	return Rule{
		Name:         "rule1",
		Terms:        []string{"rule1", "rule-1"},
		Alternatives: []string{"alt-rule1", "alt-rule-1"},
		Severity:     SevWarn,
	}
}

func TestRule_CanIgnoreLine(t *testing.T) {
	r := testRule()

	tests := []struct {
		name      string
		line      string
		assertion assert.BoolAssertionFunc
	}{
		{"violation without comment", "rule1", assert.False},
		{"violation with correct comment", "rule1 #wokeignore:rule=rule1", assert.True},
		{"violation with space as rule", "rule1 #wokeignore:rule= ", assert.False},
		{"violation with invalid comment", "rule1 #wokeignore:rule", assert.False},
		{"violation with tab as rule", "rule1 #wokeignore:rule=\t", assert.False},
		{"violation with multiple rules", "rule1 #wokeignore:rule=rule1,rule2", assert.True},
		{"violation with incorrect comment", "rule1 #wokeignore:rule=rule2", assert.False},
		{"no violation with correct comment", "rule2 #wokeignore:rule=rule1", assert.True},
		{"violation with text after ignore", "rule1 #wokeignore:rule=rule1 something else", assert.True},
		{"violation with multiple ignores", "rule1 #wokeignore:rule=rule1 wokeignore:rule=rule2", assert.True},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, r.CanIgnoreLine(tt.line))
		})
	}

}
