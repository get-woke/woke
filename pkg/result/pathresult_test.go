package result

import (
	"testing"

	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestMatchPathRules(t *testing.T) {
	tt := []struct {
		path    string
		matches int
	}{
		{path: "/foo/bar/whitelist_test.go", matches: 1},
		{path: "/foo/whitelist/blacklist_test.go", matches: 2},
		{path: "/white/list/.whitelist_test.go", matches: 1},
		{path: "/foo/bar/path_test.go", matches: 0},
		{path: "/foo/bar/whitelistblacklist-new.go", matches: 2},
	}
	for _, test := range tt {
		t.Run(test.path, func(t *testing.T) {
			pr := MatchPathRules(rule.DefaultRules, test.path)
			assert.Len(t, pr, test.matches)
			for _, p := range pr {
				assert.Equal(t, "Filename violation: "+p.Rule.Reason(p.Result.Violation), p.Reason())
			}
		})
	}
}
