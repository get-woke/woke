package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/caitlinelfring/woke/pkg/util"
	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules       []*rule.Rule `yaml:"rules,omitempty"`
	IgnoreFiles []string     `yaml:"ignore_files"`

	files []string
}

// NewConfig returns a config from the provided yaml file containing rules
func NewConfig(filename string, fileGlobs []string) (*Config, error) {
	var c Config
	var err error

	// No filename given, use default rules
	if filename == "" {
		c.AddDefaultRules()
	} else {
		err = c.load(filename)
		if len(c.Rules) == 0 {
			c.AddDefaultRules()
		}
	}

	// Ignore the config filename, it will always match on its own rules
	c.IgnoreFiles = append(c.IgnoreFiles, filename)

	allFiles, _ := util.GetFilesInGlobs(fileGlobs)
	for _, f := range allFiles {
		if c.shouldIgnoreFile(f) {
			continue
		}
		c.files = append(c.files, f)
	}

	return &c, err
}

func (c *Config) shouldIgnoreFile(f string) bool {
	for _, ignore := range c.IgnoreFiles {
		if match, _ := filepath.Match(ignore, f); match {
			return true
		}
	}
	return false
}

func (c *Config) GetFiles() []string {
	return c.files
}

func (c *Config) AddDefaultRules() {
	c.Rules = rule.DefaultRules
}

func (c *Config) load(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, c)
}
