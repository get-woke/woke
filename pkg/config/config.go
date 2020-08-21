package config

import (
	"bufio"
	"go/token"
	"io/ioutil"
	"os"

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

// Parse reads the file and returns results of places where rules are broken
func (c *Config) Parse(filename string) ([]*rule.Result, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var results []*rule.Result

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		text := scanner.Text()
		for _, r := range c.Rules {
			idx := r.Regexp.FindAllStringIndex(text, -1)
			if idx == nil {
				continue
			}

			for _, i := range idx {
				result := rule.Result{
					Rule:  r,
					Match: text[i[0]:i[1]],
					Position: &token.Position{
						Filename: filename,
						Line:     line,
						Column:   i[0],
					},
				}
				results = append(results, &result)
			}
		}

		line++
	}
	return results, scanner.Err()
}
