package config

import (
	"io/ioutil"
	"strings"

	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/caitlinelfring/woke/pkg/util"

	gitignore "github.com/sabhiram/go-gitignore"
	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules       []*rule.Rule `yaml:"rules"`
	IgnoreFiles []string     `yaml:"ignore_files"`

	_gitIgnore      *gitignore.GitIgnore
	hasAbsolutePath bool

	files []string
}

// DefaultIgnore is the default list of file globs that will be ignored
var DefaultIgnore = []string{
	".git",
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
	c.IgnoreFiles = append(c.IgnoreFiles, filename)
	c.compileIgnoreGlobs()

	return &c, nil
}

// SetFiles computes the list of files that will be checked
func (c *Config) SetFiles(fileGlobs []string) {
	allFiles, hasAbsolutePath, _ := util.GetFilesInGlobs(fileGlobs)
	c.hasAbsolutePath = hasAbsolutePath
	for _, f := range allFiles {
		if c._gitIgnore.MatchesPath(f) {
			continue
		}
		c.files = append(c.files, f)
	}
}

// compileIgnoreGlobs pre-compiles globs
// See https://github.com/gobwas/glob#performance
func (c *Config) compileIgnoreGlobs() {
	ignoreLines := []string{}
	if buffer, err := ioutil.ReadFile(".gitignore"); err == nil {
		ignoreLines = append(ignoreLines, strings.Split(string(buffer), "\n")...)
	}
	ignoreLines = append(ignoreLines, c.IgnoreFiles...)
	ignoreLines = append(ignoreLines, DefaultIgnore...)
	c._gitIgnore, _ = gitignore.CompileIgnoreLines(ignoreLines...)
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
