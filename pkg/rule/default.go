package rule

import "regexp"

// WhitelistRule is the default rule for whitelist
var WhitelistRule = Rule{
	Name:         "whitelist",
	Regexp:       regexp.MustCompile(`\b(white-?list)\b`),
	Alternatives: "allowlist",
}

// BlacklistRule is the default rule for whitelist
var BlacklistRule = Rule{
	Name:         "blacklist",
	Regexp:       regexp.MustCompile(`\b(black-?list)\b`),
	Alternatives: "denylist,blocklist",
}

// DefaultRules are the default rules in case a config file with rules is not provided
var DefaultRules = []*Rule{
	&WhitelistRule,
	&BlacklistRule,
}
