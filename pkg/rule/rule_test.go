package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRuleWithOptions(o Options) Rule {
	r := testRule()
	r.SetOptions(o)
	return r
}

func testRule() Rule {
	return Rule{
		Name:         "rule1",
		Terms:        []string{"rule1", "rule-1"},
		Alternatives: []string{"alt-rule1", "alt-rule-1"},
		Severity:     SevWarn,
	}
}

func TestRule_FindMatchIndexes(t *testing.T) {
	tests := []struct {
		text       string
		expected   [][]int
		expectedWb [][]int
	}{
		{"this string has rule-1 and rule1 included", [][]int{{16, 22}, {27, 32}}, [][]int{{16, 22}, {27, 32}}},
		{"this string has rule-2 and rule1 included", [][]int{{27, 32}}, [][]int{{27, 32}}},
		{"this string does not have any findings", [][]int(nil), [][]int(nil)},
		{"this string has finding with word boundary rule1rule-1", [][]int{{43, 48}, {48, 54}}, [][]int(nil)},
	}
	for _, test := range tests {
		r := testRule()
		got := r.FindMatchIndexes(test.text)
		assert.Equal(t, test.expected, got)
	}

	for _, test := range tests {
		r := testRuleWithOptions(Options{WordBoundary: true})
		got := r.FindMatchIndexes(test.text)
		assert.Equal(t, test.expectedWb, got)
	}

	e := Rule{Name: "rule1"}
	assert.Equal(t, [][]int(nil), e.FindMatchIndexes("rule1"))
}

func TestRule_Reason(t *testing.T) {
	r := testRule()
	assert.Equal(t, "`rule-1` may be insensitive, use `alt-rule1`, `alt-rule-1` instead", r.Reason("rule-1"))
	assert.Equal(t, "`rule1` may be insensitive, use `alt-rule1`, `alt-rule-1` instead", r.Reason(""))

	r.Alternatives = []string(nil)
	assert.Equal(t, "`rule-1` may be insensitive, try not to use it", r.Reason("rule-1"))
}

func TestRule_ReasonWithNote(t *testing.T) {
	r := testRule()

	assert.Equal(t, "`rule-1` may be insensitive, use `alt-rule1`, `alt-rule-1` instead", r.ReasonWithNote("rule-1"))

	r.Note = "rule note here"
	r.SetIncludeNote(true)
	assert.Equal(t, "`rule-1` may be insensitive, use `alt-rule1`, `alt-rule-1` instead (rule note here)", r.ReasonWithNote("rule-1"))
}

func TestRule_CanIgnoreLine(t *testing.T) {
	r := testRule()

	tests := []struct {
		name      string
		line      string
		assertion assert.BoolAssertionFunc
	}{
		{"finding without comment", "rule1", assert.False},
		{"finding with correct comment", "rule1 #wokeignore:rule=rule1", assert.True},
		{"finding with space as rule", "rule1 #wokeignore:rule= ", assert.False},
		{"finding with invalid comment", "rule1 #wokeignore:rule", assert.False},
		{"finding with tab as rule", "rule1 #wokeignore:rule=\t", assert.False},
		{"finding with multiple rules", "rule1 #wokeignore:rule=rule1,rule2", assert.True},
		{"finding with incorrect comment", "rule1 #wokeignore:rule=rule2", assert.False},
		{"no finding with correct comment", "rule2 #wokeignore:rule=rule1", assert.True},
		{"finding with text after ignore", "rule1 #wokeignore:rule=rule1 something else", assert.True},
		{"finding with multiple ignores", "rule1 #wokeignore:rule=rule1 wokeignore:rule=rule2", assert.True},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, r.CanIgnoreLine(tt.line))
		})
	}
}

func TestRule_EmptyTerms(t *testing.T) {
	r := Rule{
		Name:         "rule1",
		Terms:        []string{},
		Alternatives: []string{},
		Severity:     SevWarn,
	}
	tests := []struct {
		s         string
		wb        bool
		assertion assert.BoolAssertionFunc
	}{
		{s: "this has rule with empty terms", wb: false, assertion: assert.False},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			r.SetOptions(Options{WordBoundary: tt.wb})
			tt.assertion(t, len(r.FindMatchIndexes(tt.s)) > 0)
		})
	}
}

func TestRule_regexString(t *testing.T) {
	tests := []struct {
		desc     string
		rule     Rule
		expected string
	}{
		{
			desc:     "default",
			rule:     testRule(),
			expected: `(%s)`,
		},
		{
			desc:     "word boundary",
			rule:     testRuleWithOptions(Options{WordBoundary: true}),
			expected: `\b(%s)\b`,
		},
		{
			desc:     "word boundary start",
			rule:     testRuleWithOptions(Options{WordBoundaryStart: true}),
			expected: `\b(%s)`,
		},
		{
			desc:     "word boundary end",
			rule:     testRuleWithOptions(Options{WordBoundaryEnd: true}),
			expected: `(%s)\b`,
		},
		{
			desc:     "word boundary start and end",
			rule:     testRuleWithOptions(Options{WordBoundaryStart: true, WordBoundaryEnd: true}),
			expected: `\b(%s)\b`,
		},
		{
			// To show that enabling WordBoundary will win over other options
			desc:     "word boundary and word boundary start/end false",
			rule:     testRuleWithOptions(Options{WordBoundary: true, WordBoundaryStart: false, WordBoundaryEnd: false}),
			expected: `\b(%s)\b`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.rule.regexString())
		})
	}
}

func Test_removeInlineIgnore(t *testing.T) {
	tests := []struct {
		desc     string
		line     string
		expected string
	}{
		{
			desc:     "replace wokeignore:rule",
			line:     "wokeignore:rule=master-slave",
			expected: "����������������������������",
		},
		{
			desc:     "not replace wokeignore:rule",
			line:     "no inline ignore",
			expected: "no inline ignore",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, removeInlineIgnore(tt.line))
		})
	}
}

func TestRule_IncludeNote(t *testing.T) {
	r := testRule()
	includeNote := true

	assert.Equal(t, false, r.includeNote())

	// Test IncludeNote flag doesn't get overridden with SetIncludeNote method
	r.Options.IncludeNote = &includeNote
	r.SetIncludeNote(false)
	assert.Equal(t, true, r.includeNote())
}

func Test_IsDirectiveOnlyLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		assertion assert.BoolAssertionFunc
	}{
		{"text and no wokeignore", "some text", assert.False},
		{"text, then wokeignore", "some text #wokeignore:rule=rule1", assert.False},
		{"text, then invalid wokeignore", "some text #wokeignore:rule", assert.False},
		{"text, then multiple rules in wokeignore", "some text #wokeignore:rule=rule1,rule2", assert.False},
		{"text, then text after ignore", "some text #wokeignore:rule=rule1 something else", assert.False},
		{"text, then multiple ignores", "some text #wokeignore:rule=rule1 wokeignore:rule=rule2", assert.False},
		{"empty line", "", assert.False},
		{"only wokeignore", "#wokeignore:rule=rule1", assert.True},
		// any text to the right of wokeignore when line starts with wokeignore will not be considered by woke for findings
		{"wokeignore, then text", "#wokeignore:rule=rule1 something else", assert.True},
		{"non-alphanumeric text before and after wokeignore", "<!-- wokeignore:rule=rule1 -->", assert.True},
		{"spaces before wokeignore", "     #wokeignore:rule=rule1", assert.True},
		{"tabs before wokeignore", "\t\t\t#wokeignore:rule=rule1", assert.True},
		{"tabs and spaces before wokeignore", " \t \t \t #wokeignore:rule=rule1", assert.True},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, IsDirectiveOnlyLine(tt.line))
		})
	}
}
