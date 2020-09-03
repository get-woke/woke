package config

import (
	"bytes"
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

		_, err := NewConfig("testdata/good.yaml")
		assert.NoError(t, err)

		assert.Equal(t, `{"level":"debug","rules":["rule1","rule2","whitelist","blacklist","master-slave","slave","grandfathered","man-hours","sanity-check","dummy","he","guys"],"message":"rules enabled"}`+"\n", out.String())
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

		expectedRules = append(expectedRules, rule.DefaultRules...)

		expected := &Config{
			Rules:       expectedRules,
			IgnoreFiles: []string{"README.md", "pkg/rule/default.go", "testdata/good.yaml"},
		}

		assert.EqualValues(t, expected.Rules, c.Rules)
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

}
