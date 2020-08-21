package config

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/caitlinelfring/woke/pkg/util"
	"github.com/gobwas/glob"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

// Config contains a list of rules
type Config struct {
	Rules       []*rule.Rule `yaml:"rules"`
	IgnoreFiles []string     `yaml:"ignore_files"`

	ignoreFilesGlob []glob.Glob

	files []string
}

// DefaultIgnore is the default list of file globs that will be ignored
var DefaultIgnore = []string{
	".git/*",
}

// NewConfig returns a config from the provided yaml file containing rules
func NewConfig(filename string, fileGlobs []string) (*Config, error) {
	var c Config
	var err error

	// No filename given, use default rules
	if filename == "" {
		c.SetDefaultRules()
	} else {
		err = c.load(filename)
		if len(c.Rules) == 0 {
			c.SetDefaultRules()
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

	for _, g := range gitIgnore() {
		c.ignoreFilesGlob = append(c.ignoreFilesGlob, glob.MustCompile(g))
	}

	for _, g := range DefaultIgnore {
		c.ignoreFilesGlob = append(c.ignoreFilesGlob, glob.MustCompile(g))
	}
}

func gitIgnore() (lines []string) {
	buffer, err := os.Open(".gitignore")
	if err != nil {
		return
	}

	defer func() {
		if err = buffer.Close(); err != nil {
			log.Error().Err(err).Msg("gitignore buffer failed to close")
		}
	}()

	commentRe := regexp.MustCompile(`#(.*)$`)
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		// Remove comments
		text := commentRe.ReplaceAllString(scanner.Text(), "")
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		log.Debug().Str("entry", text).Msg("adding gitignore entry")
		lines = append(lines, text)
	}
	if err = scanner.Err(); err != nil {
		log.Info().Err(err).Msg("gitignore scanner error")
	}
	return
}

func (c *Config) ignoreFile(f string) bool {
	for _, ignore := range c.ignoreFilesGlob {
		if ignore.Match(f) {
			return true
		}
	}
	return false
}

// GetFiles returns files that may be parsed
func (c *Config) GetFiles() []string {
	return c.files
}

// SetDefaultRules sets the config Rules to DefaultRules
func (c *Config) SetDefaultRules() {
	c.Rules = rule.DefaultRules
}

func (c *Config) load(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, c)
}
