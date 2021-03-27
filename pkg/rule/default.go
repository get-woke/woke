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

//go:embed default.yaml
var defaults []byte

func init() {
	if err := yaml.Unmarshal(defaults, &DefaultRules); err != nil {
		panic(fmt.Errorf("failed to load default rules: %s", err))
	}

	for _, r := range DefaultRules {
		r.SetRegexp()
	}
}
