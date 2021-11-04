package config

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	ExcludeCategories  []string     `yaml:"exclude_categories"`
}

// NewConfig returns a new Config
func NewConfig(filename string) (*Config, error) {
	var c Config
	if len(filename) > 0 {
		var err error

		if isValidURL(filename) {
			c, err = loadRemoteConfig(filename)
		} else {
			c, err = loadConfig(filename)
		}

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

	c.ConfigureRules()
	logRuleset("all enabled", c.Rules)

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
// Filter out any rules that fall under ExcludeCategories
func (c *Config) ConfigureRules() {
	for _, r := range rule.DefaultRules {
		if !c.inExistingRules(r) {
			c.Rules = append(c.Rules, r)
		}
	}

	logRuleset("default", rule.DefaultRules)
	var excludeIndices []int

RuleLoop:
	for i, r := range c.Rules {
		for _, ex := range c.ExcludeCategories {
			// append and continue to next rule if category match found
			if r.ContainsCategory(ex) {
				excludeIndices = append(excludeIndices, i)
				continue RuleLoop
			}
		}

		r.SetRegexp()
		r.SetIncludeNote(c.IncludeNote)
	}

	// Remove excluded rules after done iterating through them
	if len(c.ExcludeCategories) > 0 {
		log.Debug().Strs("categories", c.ExcludeCategories).Msg("excluding categories")
	}
	for i, exIdx := range excludeIndices {
		// every time a rule is removed, index of rules that come after it must be reduced by one
		adjustedIdx := exIdx - i
		log.Debug().
			Strs("categories", c.Rules[adjustedIdx].Options.Categories).
			Msg(fmt.Sprintf("rule \"%s\" excluded with categories", c.Rules[adjustedIdx].Name))
		c.RemoveRule(adjustedIdx)
	}
}

// Remove rule at index i in c.Rules while maintaining order
func (c *Config) RemoveRule(i int) {
	if i >= len(c.Rules) || i < 0 {
		return
	}
	c.Rules = append(c.Rules[:i], c.Rules[i+1:]...)
}

func loadConfig(filename string) (c Config, err error) {
	yamlFile, err := ioutil.ReadFile(filename)
	log.Debug().Str("filename", filename).Msg("Adding custom ruleset from")
	if err != nil {
		return c, err
	}
	return c, yaml.Unmarshal(yamlFile, &c)
}

// gets the remote config from the url provided and returns config
func loadRemoteConfig(url string) (c Config, err error) {
	log.Debug().Str("url", url).Msg("Downloading file from")
	client := &http.Client{}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return c, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return c, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()

	// only parse response body if it is in the response is in the 2xx range
	statusOK := resp.StatusCode >= 200 && resp.StatusCode <= 299
	if !statusOK {
		return c, fmt.Errorf("unable to download remote config from url. Response code: %v. Response body: %c", resp.StatusCode, body)
	}

	log.Debug().Int("HTTP Response Status:", resp.StatusCode).Msg("Valid URL Response")
	return c, yaml.Unmarshal(body, &c)
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
		log.Debug().Strs("rules", enabledRules).Msg(fmt.Sprintf("%s rules", name))
	}
}
