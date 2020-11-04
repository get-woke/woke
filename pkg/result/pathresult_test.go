package result

import (
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/stretchr/testify/assert"
)

func TestMatchPathRules(t *testing.T) {
	testRule := rule.NewTestRule()
	secondTestRule := rule.Rule{
		Name:         "rule-2",
		Terms:        []string{"rule-2"},
		Alternatives: []string{"alt-2"},
	}
	testRules := []*rule.Rule{
		&testRule,
		&secondTestRule,
	}
	tt := []struct {
		path    string
		matches int
	}{
		{path: "/foo/bar/testrule_test.go", matches: 1},
		{path: "/foo/testrule/rule-2_test.go", matches: 2},
		{path: "/white/list/.testrule_test.go", matches: 1},
		{path: "/foo/bar/path_test.go", matches: 0},
		{path: "/foo/bar/testrule-rule-2-new.go", matches: 2},
	}
	for _, test := range tt {
		t.Run(test.path, func(t *testing.T) {
			pr := MatchPathRules(testRules, test.path)
			assert.Len(t, pr, test.matches)
			for _, p := range pr {
				assert.Equal(t, "Filename violation: "+p.Rule.Reason(p.LineResult.Violation), p.Reason())
			}
		})
	}
}
