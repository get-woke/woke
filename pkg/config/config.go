package config

import (
	"io/ioutil"
	"os"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var defaultConfigFilenames = []string{".woke.yaml", ".woke.yml"}

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
	} else if defaultCfg := loadDefaultConfigFiles(); defaultCfg != nil {
		c = *defaultCfg

		c.IgnoreFiles = append(c.IgnoreFiles, defaultConfigFilenames...)
	}

	c.ConfigureRules()

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

// ConfigureRules adds the config Rules to DefaultRules
func (c *Config) ConfigureRules() {
	for _, r := range rule.DefaultRules {
		if !c.inExistingRules(r) {
			c.Rules = append(c.Rules, r)
		}
	}

	for _, r := range c.Rules {
		r.SetRegexp()
	}
}

func (c *Config) load(filename string) error {
	cfg, err := loadConfig(filename)
	if err != nil {
		return err
	}
	*c = *cfg
	return nil
}

func loadConfig(filename string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var c Config
	err = yaml.Unmarshal(yamlFile, &c)

	return &c, err
}

func loadDefaultConfigFiles() (cfg *Config) {
	for _, file := range defaultConfigFilenames {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Debug().Str("cfg", file).Err(err).Msg("tried default config file")
			continue
		}
		var err error
		cfg, err = loadConfig(file)
		if err == nil && cfg != nil {
			log.Debug().Str("cfg", file).Msg("found default config file!")
			return
		}
	}
	return
}
