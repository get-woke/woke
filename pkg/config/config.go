package config

import (
	"io/ioutil"

	"github.com/get-woke/woke/pkg/rule"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules       []*rule.Rule `yaml:"rules"`
	IgnoreFiles []string     `yaml:"ignore_files"`
}

// NewConfig returns a new Config
func NewConfig(filename string) (*Config, error) {
	var c Config

	if filename != "" {
		if err := c.load(filename); err != nil {
			return nil, err
		}
		// Ignore the config filename, it will always match on its own rules
		c.IgnoreFiles = append(c.IgnoreFiles, filename)
	}

	c.AddDefaultRules()

	// For debugging/informational purposes
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		enabledRules := make([]string, len(c.Rules))
		for i := range c.Rules {
			enabledRules[i] = c.Rules[i].Name
		}
		log.Debug().Strs("rules", enabledRules).Msg("rules enabled")
	}

	return &c, nil
}

func (c *Config) inExistingRules(r *rule.Rule) bool {
	for _, n := range c.Rules {
		if n.Name == r.Name {
			return true
		}
	}
	return false
}

// AddDefaultRules adds the config Rules to DefaultRules
func (c *Config) AddDefaultRules() {
	for _, r := range rule.DefaultRules {
		if !c.inExistingRules(r) {
			c.Rules = append(c.Rules, r)
		}
	}
}

func (c *Config) load(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, c)
}
