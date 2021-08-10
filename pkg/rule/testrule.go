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

var TestErrorRule = Rule{
	Name:         "slave",           // wokeignore:rule=slave
	Terms:        []string{"slave"}, // wokeignore:rule=slave
	Alternatives: []string{"follower"},
	Severity:     0,
	Options: Options{
		WordBoundary: false,
	},
}

var TestInfoRule = Rule{
	Name:         "test",
	Terms:        []string{"test"},
	Alternatives: []string{"alternative"},
	Severity:     2,
	Options: Options{
		WordBoundary: false,
	},
}

func init() {
	TestRule.SetRegexp()
	TestErrorRule.SetRegexp()
	TestInfoRule.SetRegexp()
}
