package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules              []*rule.Rule `yaml:"rules"`
	IgnoreFiles        []string     `yaml:"ignore_files"`
	SuccessExitMessage *string      `yaml:"success_exit_message"`
	IncludeNote        bool         `yaml:"include_note"`
}

// NewConfig returns a new Config
func NewConfig(filename string, disableDefaultRules bool) (*Config, error) {
	var c Config
	if len(filename) > 0 {
		var err error
		c, err = loadConfig(filename)
		if err != nil {
			return nil, err
		}

		log.Debug().Str("config", filename).Msg("loaded config file")
		logRuleset("config", c.Rules)

		// Ignore the config filename, it will always match on its own rules
		c.IgnoreFiles = append(c.IgnoreFiles, relative(filename))
	} else {
		log.Debug().Msg("no config file loaded, using only default rules")
	}

	c.ConfigureRules(disableDefaultRules)
	logRuleset("all", c.Rules)

	return &c, nil
}

// GetSuccessExitMessage returns the message to be shows on a successful exit as
// defined in the config, or a default message.
func (c *Config) GetSuccessExitMessage() string {
	if c.SuccessExitMessage == nil {
		return "No findings found."
	}
	return *c.SuccessExitMessage
}

func (c *Config) inExistingRules(r *rule.Rule) bool {
	for _, n := range c.Rules {
		if n.Name == r.Name {
			return true
		}
	}
	return false
}

// ConfigureRules adds the config Rules to DefaultRules
// Configure RegExps for all rules
// Configure IncludeNote for all rules
func (c *Config) ConfigureRules(disableDefaultRules bool) {
	if !disableDefaultRules {
		for _, r := range rule.DefaultRules {
			if !c.inExistingRules(r) {
				c.Rules = append(c.Rules, r)
			}
		}

		logRuleset("default", rule.DefaultRules)
	}

	for _, r := range c.Rules {
		r.SetRegexp()
		r.SetIncludeNote(c.IncludeNote)
	}
}

func loadConfig(filename string) (c Config, err error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}

	return c, yaml.Unmarshal(yamlFile, &c)
}

func relative(filename string) string {
	// viper provides an absolute path to the config file, but we want the relative
	// path to the config file from the current directory to make it easy for woke to ignore it
	if filepath.IsAbs(filename) {
		cwd, _ := os.Getwd()
		if relfilename, err := filepath.Rel(cwd, filename); err == nil {
			return relfilename
		}
	}
	return filename
}

// For debugging/informational purposes
func logRuleset(name string, rules []*rule.Rule) {
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		enabledRules := make([]string, len(rules))
		for i := range rules {
			enabledRules[i] = rules[i].Name
		}
		log.Debug().Strs("rules", enabledRules).Msg(fmt.Sprintf("%s rules enabled", name))
	}
}
