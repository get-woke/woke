// Code generated by internal/rule/gen.go; DO NOT EDIT.
// Regenerate with 'go generate ./...'

package rule

// BlackboxRule is the default rule for "blackbox"
var BlackboxRule = Rule{
	Name:         "blackbox",
	Terms:        []string{"black-box", "blackbox", "black box"},
	Alternatives: []string{"closed-box"},
	Options: Options{
		WordBoundary: false,
	},
}

// BlacklistRule is the default rule for "blacklist"
var BlacklistRule = Rule{
	Name:         "blacklist",
	Terms:        []string{"blacklist", "black-list", "blacklisted", "black-listed"},
	Alternatives: []string{"denylist", "blocklist"},
	Severity:     1,
	Options: Options{
		WordBoundary: false,
	},
}

// DummyRule is the default rule for "dummy"
var DummyRule = Rule{
	Name:         "dummy",
	Terms:        []string{"dummy"},
	Alternatives: []string{"placeholder", "sample"},
	Options: Options{
		WordBoundary: false,
	},
}

// GrandfatheredRule is the default rule for "grandfathered"
var GrandfatheredRule = Rule{
	Name:         "grandfathered",
	Terms:        []string{"grandfathered"},
	Alternatives: []string{"legacy status"},
	Options: Options{
		WordBoundary: false,
	},
}

// GuysRule is the default rule for "guys"
var GuysRule = Rule{
	Name:         "guys",
	Terms:        []string{"guys"},
	Alternatives: []string{"everyone", "folks", "people", "you all", "y'all", "yinz"},
	Options: Options{
		WordBoundary: false,
	},
}

var HeSheRule = Rule{
	Name:  "he-she-rule",
	Terms: []string{"he", "she", "hers", "his", "he'd", "she'd", "she'll", "he'll", "she's", "he's"},
	Alternatives: []string{"they", "them", "their", "theirs", "it"},
	Options: Options{
		WordBoundary: false,
	},
}

// ManHoursRule is the default rule for "man-hours"
var ManHoursRule = Rule{
	Name:         "man-hours",
	Terms:        []string{"man hours", "man-hours"},
	Alternatives: []string{"person hours", "engineer hours"},
	Options: Options{
		WordBoundary: false,
	},
}

// MasterSlaveRule is the default rule for "master-slave"
var MasterSlaveRule = Rule{
	Name:         "master-slave",
	Terms:        []string{"master-slave", "master/slave"},
	Alternatives: []string{"leader/follower", "primary/replica", "primary/standby", "primary/secondary"},
	Options: Options{
		WordBoundary: false,
	},
}

// SanityRule is the default rule for "sanity"
var SanityRule = Rule{
	Name:         "sanity",
	Terms:        []string{"sanity"},
	Alternatives: []string{"confidence", "quick check", "coherence check"},
	Options: Options{
		WordBoundary: false,
	},
}

// SlaveRule is the default rule for "slave"
var SlaveRule = Rule{
	Name:         "slave",
	Terms:        []string{"slave"},
	Alternatives: []string{"follower", "replica", "standby"},
	Options: Options{
		WordBoundary: false,
	},
}

// WhiteboxRule is the default rule for "whitebox"
var WhiteboxRule = Rule{
	Name:         "whitebox",
	Terms:        []string{"white-box", "whitebox", "white box"},
	Alternatives: []string{"open-box"},
	Options: Options{
		WordBoundary: false,
	},
}

// WhitelistRule is the default rule for "whitelist"
var WhitelistRule = Rule{
	Name:         "whitelist",
	Terms:        []string{"whitelist", "white-list", "whitelisted", "white-listed"},
	Alternatives: []string{"allowlist"},
	Severity:     1,
	Options: Options{
		WordBoundary: false,
	},
}

// DefaultRules are the default rules always used
var DefaultRules = []*Rule{
	&BlackboxRule,
	&BlacklistRule,
	&DummyRule,
	&GrandfatheredRule,
	&GuysRule,
	&HeSheRule,
	&ManHoursRule,
	&MasterSlaveRule,
	&SanityRule,
	&SlaveRule,
	&WhiteboxRule,
	&WhitelistRule,
}
