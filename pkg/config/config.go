package config

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
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

func NewConfig(filename string) (*Config, error) {
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

	return &c, err
}

// SetFiles computes the list of files that will be checked
func (c *Config) SetFiles(fileGlobs []string) {
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
		c.addIgnoreGlob(ignore)
	}

	for _, ignore := range gitIgnore() {
		c.addIgnoreGlob(ignore)
	}

	for _, ignore := range DefaultIgnore {
		c.addIgnoreGlob(ignore)
	}
}

func (c *Config) addIgnoreGlob(s string) {
	abs, err := filepath.Abs(s)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", s).
			Msg("failed to get absolute path")
		return
	}
	c.ignoreFilesGlob = append(c.ignoreFilesGlob, glob.MustCompile(abs))
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
