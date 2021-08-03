package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/get-woke/woke/pkg/rule"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("check-logger", func(t *testing.T) {
		out := &bytes.Buffer{}
		log.Logger = zerolog.New(out)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		// isolate config rules
		loaded, err := loadConfig("testdata/good.yaml")
		assert.NoError(t, err)
		configRules := make([]string, len(loaded.Rules))
		for i := range loaded.Rules {
			configRules[i] = fmt.Sprintf("%q", loaded.Rules[i].Name)
		}

		// isolate default rules
		defaultRules := make([]string, len(rule.DefaultRules))
		for i := range defaultRules {
			defaultRules[i] = fmt.Sprintf("%q", rule.DefaultRules[i].Name)
		}

		c, err := NewConfig("testdata/good.yaml")
		assert.NoError(t, err)
		enabledRules := make([]string, len(c.Rules))
		for i := range c.Rules {
			enabledRules[i] = fmt.Sprintf("%q", c.Rules[i].Name)
		}

		loadedRemoteConfigMsg := `{"level":"debug","filename":"testdata/good.yaml","message":"Adding custom ruleset from"}`
		loadedRemoteConfig := `{"level":"debug","filename":"testdata/good.yaml","message":"Adding custom ruleset from"}`
		loadedConfigMsg := `{"level":"debug","config":"testdata/good.yaml","message":"loaded config file"}`
		configRulesMsg := fmt.Sprintf(`{"level":"debug","rules":[%s],"message":"config rules enabled"}`, strings.Join(configRules, ","))
		defaultRulesMsg := fmt.Sprintf(`{"level":"debug","rules":[%s],"message":"default rules enabled"}`, strings.Join(defaultRules, ","))
		allRulesMsg := fmt.Sprintf(`{"level":"debug","rules":[%s],"message":"all rules enabled"}`, strings.Join(enabledRules, ","))
		assert.Equal(t,
			loadedRemoteConfigMsg+"\n"+loadedRemoteConfig+"\n"+loadedConfigMsg+"\n"+configRulesMsg+"\n"+defaultRulesMsg+"\n"+allRulesMsg+"\n",
			out.String())
	})

	t.Run("config-good", func(t *testing.T) {
		c, err := NewConfig("testdata/good.yaml")
		assert.NoError(t, err)

		expectedRules := []*rule.Rule{}
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule1",
			Terms:        []string{"rule1"},
			Alternatives: []string{"alt-rule1"},
			Severity:     rule.SevWarn,
		})
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "rule2",
			Terms:        []string{"rule2", "rule-2"},
			Alternatives: []string{"alt-rule2", "alt-rule-2"},
			Severity:     rule.SevError,
		})
		expectedRules = append(expectedRules, &rule.Rule{
			Name:         "whitelist",
			Terms:        []string{"rulewl", "rule-wl"},
			Alternatives: []string{"alt-rulewl", "alt-rule-wl"},
			Severity:     rule.SevError,
		})

		expected := &Config{
			Rules:       expectedRules,
			IgnoreFiles: []string{"README.md", "pkg/rule/default.go", "testdata/good.yaml"},
		}
		expected.ConfigureRules()

		assert.EqualValues(t, expected.Rules, c.Rules)

		// check default config message
		assert.Equal(t, "No findings found.", c.GetSuccessExitMessage())
	})

	t.Run("config-empty-missing", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("")
		assert.NoError(t, err)

		expectedEmpty := &Config{
			Rules:       rule.DefaultRules,
			IgnoreFiles: []string(nil),
		}
		assert.Equal(t, expectedEmpty, c)
	})

	t.Run("config-missing", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("testdata/missing.yaml")
		assert.Error(t, err)
		assert.Nil(t, c)
	})

	t.Run("config-empty-success-message", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("testdata/empty-success-message.yaml")
		assert.NoError(t, err)

		// check default config message
		assert.Equal(t, "", c.GetSuccessExitMessage())
	})

	t.Run("config-custom-success-message", func(t *testing.T) {
		// Test when no config file is provided
		c, err := NewConfig("testdata/custom-success-message.yaml")
		assert.NoError(t, err)

		// check default config message
		assert.Equal(t, "this is a test", c.GetSuccessExitMessage())
	})

	t.Run("config-add-note-messaage", func(t *testing.T) {
		// Test when it is configured to add a note to the output message
		c, err := NewConfig("testdata/add-note-message.yaml")
		assert.NoError(t, err)

		// check global IncludeNote
		assert.Equal(t, true, c.IncludeNote)

		// check IncludeNote is set for rule2
		assert.Equal(t, true, *c.Rules[1].Options.IncludeNote)

		// check IncludeNote is not overridden for rule1
		assert.Equal(t, false, *c.Rules[0].Options.IncludeNote)
	})

	t.Run("config-dont-add-note-messaage", func(t *testing.T) {
		// Test when it is nott configured to add a note to the output message
		c, err := NewConfig("testdata/dont-add-note-message.yaml")
		assert.NoError(t, err)

		// check global IncludeNote
		assert.Equal(t, false, c.IncludeNote)

		// check IncludeNote is not set for rule2
		assert.Equal(t, false, *c.Rules[1].Options.IncludeNote)

		// check IncludeNote is not overridden for rule1
		assert.Equal(t, true, *c.Rules[0].Options.IncludeNote)
	})
}

func Test_LoadConfig(t *testing.T) {
	t.Run("valid-url", func(t *testing.T) {
		c, err := loadConfig("https://raw.githubusercontent.com/get-woke/woke/main/example.yaml")
		assert.NoError(t, err)
		assert.NotNil(t, c)
		//delete downloaded file after successful test
		os.Remove("downloadedRules.yaml")
	})

	t.Run("invalid-url", func(t *testing.T) {})
	_, err := loadConfig("https://raw.githubusercontent.com/get-woke/woke/main/example")
	assert.Error(t, err)
}

func Test_isValidURL(t *testing.T) {
	t.Run("valid-url", func(t *testing.T) {
		boolResponse := isValidUrl("https://raw.githubusercontent.com/get-woke/woke/main/example.yaml")
		assert.True(t, boolResponse)
	})

	t.Run("invalid-url", func(t *testing.T) {
		boolResponse := isValidUrl("Users/Document/test.yaml")
		assert.False(t, boolResponse)
	})

	t.Run("invalid-url", func(t *testing.T) {
		boolResponse := isValidUrl("/Users/Document/test.yaml")
		assert.False(t, boolResponse)
	})

}

func Test_DownloadFile(t *testing.T) {
	t.Run("valid-url", func(t *testing.T) {
		err := DownloadFile("DownloadedFile.yaml", "https://raw.githubusercontent.com/get-woke/woke/main/example.yaml")
		assert.NoError(t, err)
		//delete downloaded file after successful test
		os.Remove("DownloadedFile.yaml")
	})

	t.Run("invalid-url", func(t *testing.T) {
		err := DownloadFile("DownloadedFile.yaml", "https://raw.githubusercontent.com/get-woke/woke/main/example")
		assert.Error(t, err)
	})

}

func Test_relative(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	assert.Equal(t, ".woke.yml", relative(filepath.Join(cwd, ".woke.yml")))
	assert.Equal(t, ".woke.yml", relative(".woke.yml"))
	assert.Equal(t, "dir/.woke.yml", relative("dir/.woke.yml"))
}
