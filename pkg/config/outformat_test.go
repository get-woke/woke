package config

import (
	"testing"

	"github.com/get-woke/woke/pkg/printer"
	"github.com/stretchr/testify/assert"
)

func TestCreatePrinter(t *testing.T) {
	text, err := CreatePrinter(OutFormatText)
	assert.NoError(t, err)
	assert.IsType(t, &printer.Text{}, text)

	simple, err := CreatePrinter(OutFormatSimple)
	assert.NoError(t, err)
	assert.IsType(t, &printer.Simple{}, simple)

	gha, err := CreatePrinter(OutFormatGitHubActions)
	assert.NoError(t, err)
	assert.IsType(t, &printer.GitHubActions{}, gha)

	_, err = CreatePrinter("invalid-printer")
	assert.Errorf(t, err, "invalid-printer is not a valid printer type")
}
