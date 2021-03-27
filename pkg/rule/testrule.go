package rule

// TestRule is only meant to be used in tests.
// TODO: Use test terms
var TestRule = Rule{
	Name:         "whitelist",                                                        // wokeignore:rule=whitelist
	Terms:        []string{"whitelist", "white-list", "whitelisted", "white-listed"}, // wokeignore:rule=whitelist
	Alternatives: []string{"allowlist"},
	Severity:     1,
	Options: Options{
		WordBoundary: false,
	},
}

func init() {
	TestRule.SetRegexp()
}
