package rule

// WhitelistRule is the default rule for whitelist
var WhitelistRule = Rule{
	Name:         "whitelist",
	Terms:        []string{"whitelist", "white-list"},
	Alternatives: []string{"allowlist"},
	Severity:     SevWarn,
}

// BlacklistRule is the default rule for blacklist
var BlacklistRule = Rule{
	Name:         "blacklist",
	Terms:        []string{"blacklist", "black-list"},
	Alternatives: []string{"blocklist"},
	Severity:     SevWarn,
}

// DefaultRules are the default rules always used
var DefaultRules = []*Rule{
	&WhitelistRule,
	&BlacklistRule,
}
