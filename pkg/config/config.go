package config

import (
	"io/ioutil"

	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/caitlinelfring/woke/pkg/util"
	"github.com/gobwas/glob"

	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules       []*rule.Rule `yaml:"rules"`
	IgnoreFiles []string     `yaml:"ignore_files"`

	ignoreFilesGlob []glob.Glob

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
	c.compileIgnoreGlobs()

	// Must come after compiling ignore globs
	c.setFiles(fileGlobs)

	return &c, err
}

func (c *Config) setFiles(fileGlobs []string) {
	allFiles, _ := util.GetFilesInGlobs(fileGlobs)
	for _, f := range allFiles {
		if c.ignoreFile(f) {
			continue
		}
		c.files = append(c.files, f)
	}
}

// compileIgnoreGlobs pre-compiles globs
// See https://github.com/gobwas/glob#performance
func (c *Config) compileIgnoreGlobs() {
	c.ignoreFilesGlob = make([]glob.Glob, 0)

	for _, ignore := range c.IgnoreFiles {
		c.ignoreFilesGlob = append(c.ignoreFilesGlob, glob.MustCompile(ignore))
	}
}

func (c *Config) ignoreFile(f string) bool {
	for _, ignore := range c.ignoreFilesGlob {
		if ignore.Match(f) {
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
