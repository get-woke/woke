package rule

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRule_FindMatchIndexes(t *testing.T) {
	tests := []struct {
		text       string
		expected   [][]int
		expectedWb [][]int
	}{
		{"this string has test-rule and testrule included", [][]int{{16, 25}, {30, 38}}, [][]int{{16, 25}, {30, 38}}},
		{"this string has no-rule and testrule included", [][]int{{28, 36}}, [][]int{{28, 36}}},
		{"this string does not have any violations", [][]int(nil), [][]int(nil)},
		{"this string has violation with word boundary testruletest-rule", [][]int{{45, 53}, {53, 62}}, [][]int(nil)},
	}
	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			r := NewTestRule()
			got := r.FindMatchIndexes(test.text)
			assert.Equal(t, test.expected, got)
		})
	}

	for _, test := range tests {
		t.Run("word_boundary_"+test.text, func(t *testing.T) {
			r := NewTestRule()
			r.Options.WordBoundary = true

			got := r.FindMatchIndexes(test.text)
			assert.Equal(t, test.expectedWb, got)
		})
	}

	e := Rule{Name: "rule1"}
	assert.Equal(t, [][]int(nil), e.FindMatchIndexes("rule1"))
}

func TestRule_Reason(t *testing.T) {
	r := NewTestRule()
	assert.Equal(t, "`test-rule` may be insensitive, use `better-rule` instead", r.Reason("test-rule"))

	assert.Equal(t, "`test-rule` may be insensitive, use `better-rule` instead", r.Reason("test-rule"))
	assert.Equal(t, "`test-rule` may be insensitive, use `better-rule` instead", r.Reason(""))

	r.Alternatives = []string(nil)
	assert.Equal(t, "`test-rule` may be insensitive, try not to use it", r.Reason("test-rule"))
}

func TestRule_ReasonWithNote(t *testing.T) {
	r := NewTestRule()

	assert.Equal(t, "`test-rule` may be insensitive, use `better-rule` instead", r.ReasonWithNote("test-rule"))

	r.Note = "rule note here"
	assert.Equal(t, "`test-rule` may be insensitive, use `better-rule` instead (rule note here)", r.ReasonWithNote("test-rule"))
}

func TestRule_CanIgnoreLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		assertion assert.BoolAssertionFunc
	}{
		{"violation without comment", "test-rule", assert.False},
		{"violation with correct comment", "test-rule #wokeignore:rule=test-rule", assert.True},
		{"violation with space as rule", "test-rule #wokeignore:rule= ", assert.False},
		{"violation with invalid comment", "test-rule #wokeignore:rule", assert.False},
		{"violation with tab as rule", "test-rule #wokeignore:rule=\t", assert.False},
		{"violation with multiple rules", "test-rule #wokeignore:rule=test-rule,rule2", assert.True},
		{"violation with incorrect comment", "test-rule #wokeignore:rule=rule2", assert.False},
		{"no violation with correct comment", "rule2 #wokeignore:rule=test-rule", assert.True},
		{"violation with text after ignore", "test-rule #wokeignore:rule=test-rule something else", assert.True},
		{"violation with multiple ignores", "test-rule #wokeignore:rule=test-rule wokeignore:rule=rule2", assert.True},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewTestRule()
			tt.assertion(t, r.CanIgnoreLine(tt.line))
		})
	}
}

func TestRule_MatchString(t *testing.T) {
	tests := []struct {
		s         string
		wb        bool
		assertion assert.BoolAssertionFunc
	}{
		{s: "this has testrule in the middle with word boundaries", wb: true, assertion: assert.True},
		{s: "this has testrule in the middle", wb: false, assertion: assert.True},
		{s: "testruleshouldn't match with word boundaries", wb: true, assertion: assert.False},
		{s: "testruleshould match without word boundaries", wb: false, assertion: assert.True},
		{s: "thistestruleshould match without word boundaries", wb: false, assertion: assert.True},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			r := NewTestRule()
			fmt.Println(r.MatchString(tt.s, tt.wb), tt.s)
			tt.assertion(t, r.MatchString(tt.s, tt.wb))
		})
	}
}

func TestRule_EmptyTerms(t *testing.T) {
	tests := []struct {
		s         string
		wb        bool
		assertion assert.BoolAssertionFunc
	}{
		{s: "this has rule with empty terms", wb: false, assertion: assert.False},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			r := NewTestRule()
			r.Terms = []string{}

			fmt.Println(r.MatchString(tt.s, tt.wb), tt.s)
			tt.assertion(t, r.MatchString(tt.s, tt.wb))
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
