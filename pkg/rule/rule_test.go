package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRule_FindAllStringIndex(t *testing.T) {
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
		got := r.FindAllStringIndex(test.text)
		assert.Equal(t, test.expected, got)
	}
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
