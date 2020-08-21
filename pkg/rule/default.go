package rule

import "regexp"

var WhitelistRule = Rule{
	Name:         "whitelist",
	Regexp:       regexp.MustCompile(`\b(white-?list)\b`),
	Alternatives: "allowlist",
}

var BlacklistRule = Rule{
	Name:         "blacklist",
	Regexp:       regexp.MustCompile(`\b(black-?list)\b`),
	Alternatives: "denylist,blocklist",
}

var DefaultRules = []*Rule{&WhitelistRule, &BlacklistRule}
