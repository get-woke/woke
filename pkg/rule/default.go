package rule

import (
	// empty import required as a part of the embed package
	// https://golang.google.cn/pkg/embed/
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v2"
)

// DefaultRules are the default rules always used.
// This will be populated by the embed package on init
var DefaultRules = []*Rule{}

// WhitelistRule is the default rule for "whitelist" # wokeignore:rule=whitelist
// FIXME: these are only used by tests right now. The tests need to be refactored so they can be removed
var WhitelistRule *Rule // wokeignore:rule=whitelist

//go:embed default.yaml
var defaults []byte

func init() {
	if err := loadDefaultRules(); err != nil {
		panic(fmt.Errorf("failed to load default rules: %s", err))
	}

	WhitelistRule = getDefaultRule("whitelist") // wokeignore:rule=whitelist
}

func loadDefaultRules() error {
	if err := yaml.Unmarshal(defaults, &DefaultRules); err != nil {
		return err
	}
	for _, r := range DefaultRules {
		r.SetRegexp()
	}
	return nil
}

func getDefaultRule(name string) *Rule {
	if len(DefaultRules) == 0 {
		_ = loadDefaultRules()
	}
	for _, r := range DefaultRules {
		if r.Name == name {
			return r
		}
	}
	return nil
}
