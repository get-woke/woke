package config

import (
	"io/ioutil"

	"github.com/caitlinelfring/woke/pkg/rule"
	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules []*rule.Rule `yaml:"rules"`
}

// NewConfig returns a config from the provided yaml file containing rules
func NewConfig(filename string) (*Config, error) {
	var c Config
	err := c.load(filename)
	return &c, err
}

func (c *Config) load(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, c)
}
