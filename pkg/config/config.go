package config

import (
	"io/ioutil"

	"github.com/caitlinelfring/woke/pkg/ignore"
	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/caitlinelfring/woke/pkg/util"

	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules       []*rule.Rule `yaml:"rules"`
	IgnoreFiles []string     `yaml:"ignore_files"`

	hasAbsolutePath bool

	ignoreMatcherFunc func(string) bool

	files []string
}

func NewConfig(filename string) (*Config, error) {
	var c Config

	if filename != "" {
		if err := c.load(filename); err != nil {
			return &c, err
		}
	}
	c.AddDefaultRules()

	// Ignore the config filename, it will always match on its own rules
	ignorer, err := ignore.NewIgnore(append(c.IgnoreFiles, filename))
	if err != nil {
		return &c, err
	}
	c.ignoreMatcherFunc = ignorer.Match

	return &c, nil
}

// SetFiles computes the list of files that will be checked
func (c *Config) SetFiles(fileGlobs []string) {
	allFiles, hasAbsolutePath, _ := util.GetFilesInGlobs(fileGlobs)
	c.hasAbsolutePath = hasAbsolutePath
	for _, f := range allFiles {
		if c.ignoreMatcherFunc(f) {
			continue
		}
		c.files = append(c.files, f)
	}
}

// GetFiles returns files that may be parsed
func (c *Config) GetFiles() []string {
	return c.files
}

// AddDefaultRules adds the config Rules to DefaultRules
func (c *Config) AddDefaultRules() {
	c.Rules = append(c.Rules, rule.DefaultRules...)
}

func (c *Config) load(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, c)
}
