package rule

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestSeverity_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		input    string
		expected Severity
	}{
		{"warn", SevWarn},
		{"warning", SevWarn},
		{"error", SevError},
		{"info", SevInfo},
		{"not-valid", SevInfo},
	}
	for _, test := range tests {
		sev := new(Severity)
		err := yaml.Unmarshal([]byte(test.input), &sev)
		assert.NoError(t, err)

		assert.Equalf(t, test.expected, *sev, "expected: %s, got: %s", test.expected, sev)
	}
}

func TestSeverity_Colorize(t *testing.T) {
	tests := []struct {
		input    Severity
		expected string
	}{
		{SevWarn, "\x1b[33mwarning\x1b[0m"},
		{SevError, "\x1b[31merror\x1b[0m"},
		{SevInfo, "\x1b[32minfo\x1b[0m"},
		{Severity(999), "\x1b[32minfo\x1b[0m"},
	}

	color.NoColor = false
	for _, test := range tests {
		assert.Equalf(t, test.expected, test.input.Colorize(), "severity: %s", test.input)
	}
}
