package config

import (
	"testing"

	"github.com/get-woke/woke/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
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

	expected := &Config{
		Rules:       expectedRules,
		IgnoreFiles: []string{"README.md", "pkg/rule/default.go", "testdata/good.yaml"},
	}

	assert.EqualValues(t, expected, c)

	// Test when no config file is provided
	c, err = NewConfig("")
	assert.NoError(t, err)

	expectedEmpty := &Config{
		Rules:       rule.DefaultRules,
		IgnoreFiles: []string(nil),
	}
	assert.Equal(t, expectedEmpty, c)

	missing, err := NewConfig("testdata/missing.yaml")
	assert.Error(t, err)
	assert.Equal(t, expectedEmpty, missing)
}
