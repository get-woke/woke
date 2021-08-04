package config

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

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
func NewConfig(filename string) (*Config, error) {
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

	c.ConfigureRules()
	logRuleset("all", c.Rules)

	return &c, nil
}

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
func (c *Config) ConfigureRules() {
	for _, r := range rule.DefaultRules {
		if !c.inExistingRules(r) {
			c.Rules = append(c.Rules, r)
		}
	}

	logRuleset("default", rule.DefaultRules)

	for _, r := range c.Rules {
		r.SetRegexp()
		r.SetIncludeNote(c.IncludeNote)
	}
}

func loadConfig(filename string) (c Config, err error) {
	if isValidURL(filename) {
		// if fileName is a valid URL, we will download and set it to the config
		log.Debug().Str("url", filename).Msg("Downloading file from")
		// hardcoding this file and saving to root directory
		downloadedFile := "downloadedRules.yaml"
		err := DownloadFile(downloadedFile, filename)
		if err != nil {
			return c, err
		}
		filename = downloadedFile
		log.Debug().Str("filename", filename).Msg("Saved remote config to local file.")
	}

	yamlFile, err := ioutil.ReadFile(filename)
	log.Debug().Str("filename", filename).Msg("Adding custom ruleset from")
	if err != nil {
		return c, err
	}
	return c, yaml.Unmarshal(yamlFile, &c)
}

// isValidUrl tests a string to determine if it is a valid URL or not
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	log.Debug().Str("remoteConfig", toTest).Msg("Valid URL for remote config.")
	return true
}

// downloads file from url to set filepath
func DownloadFile(filepath string, url string) error {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	ctx := context.Background()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// only parse response body if it is in the response is in the 2xx range
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Debug().Int("HTTP Response Status:", resp.StatusCode).Msg("Valid URL Response")
		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()
		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		return err
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unable to download remote config from url. Response code: %v. Response body: %c", resp.StatusCode, body)
	}
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
