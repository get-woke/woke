package printer

import (
	"testing"

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
	}

	for _, test := range tests {
		p, err := NewPrinter(test.OutFormat)
		assert.NoError(t, err)
		assert.IsType(t, test.Type, p)
	}

	_, err := NewPrinter("invalid-printer")
	assert.Errorf(t, err, "invalid-printer is not a valid printer type")
}
