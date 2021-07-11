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
				assert.Equal(t, "Filename finding: "+p.Rule.Reason(p.LineResult.Finding), p.Reason())
			}
		})
	}
}

func TestMatchPathRulesBoundary(t *testing.T) {
	tt := []struct {
		path    string
		matches int
	}{
		{path: "/whitelist/list/test.go", matches: 1},
		{path: "/testwhitelist/list/test.go", matches: 0},
		{path: "/whitetest/list/test.go", matches: 0},
		{path: "/whitelist.test/list/test.go", matches: 1},
		{path: "/whitelist.test/blacklist/test.go", matches: 2},
		{path: "/whitelisttest/blacklist/test.go", matches: 1},
		{path: "/foo/bar/path_test.go", matches: 0},
		{path: "/foo/bar/whitelistblacklist-new.go", matches: 0},
		{path: "/foo/bar/whitelist-new.go", matches: 1},
	}
	for _, test := range tt {
		t.Run(test.path, func(t *testing.T) {
			defaultRules := rule.DefaultRules
			for i := range defaultRules {
				defaultRules[i].SetOptions(rule.Options{WordBoundary: true})
			}
			pr := MatchPathRules(defaultRules, test.path)
			assert.Len(t, pr, test.matches)
			for _, p := range pr {
				assert.Equal(t, "Filename finding: "+p.Rule.Reason(p.LineResult.Finding), p.Reason())
			}
		})
	}
}
