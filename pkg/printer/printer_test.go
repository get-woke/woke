package printer

import (
	"io"
	"testing"

	config "github.com/get-woke/woke/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCreatePrinter(t *testing.T) {
	tests := []struct {
		OutFormat string
		Type      Printer
	}{
		{OutFormatSimple, &Simple{}},
		{OutFormatText, &Text{}},
		{OutFormatGitHubActions, &GitHubActions{}},
		{OutFormatJSON, &JSON{}},
		{OutFormatSonarQube, &SonarQube{}},
		{OutFormatCheckstyle, &Checkstyle{}},
	}

	for _, test := range tests {
		p, err := NewPrinter(test.OutFormat, io.Discard, new(config.Config))
		assert.NoError(t, err)
		assert.IsType(t, test.Type, p)
	}

	_, err := NewPrinter("invalid-printer", io.Discard, new(config.Config))
	assert.Errorf(t, err, "invalid-printer is not a valid printer type")
}
