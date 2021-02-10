package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/util"

	"github.com/rs/zerolog/log"
)

// Printer is an interface for printing FileResults
type Printer interface {
	Print(io.Writer, *result.FileResults) error
}

const (
	// OutFormatText is a text-based output format, best for CLIs
	OutFormatText = "text"
	// OutFormatSimple is a simplified output format, which can be used with something like https://github.com/reviewdog/reviewdog
	OutFormatSimple = "simple"
	// OutFormatGitHubActions is an output format supported by GitHub Actions annotations
	// https://docs.github.com/en/free-pro-team@latest/actions/reference/workflow-commands-for-github-actions#setting-a-warning-message
	OutFormatGitHubActions = "github-actions"

	// OutFormatJSON outputs in json
	OutFormatJSON = "json"
)

// OutFormats are all the available output formats. The first one should be the default
var OutFormats = []string{
	OutFormatText,
	OutFormatSimple,
	OutFormatGitHubActions,
	OutFormatJSON,
}

// OutFormatsString is all OutFormats, as a comma-separated string
var OutFormatsString = strings.Join(OutFormats, ",")

// NewPrinter returns a valid new Printer from a string, or an error if the printer is invalid
func NewPrinter(f string) (Printer, error) {
	var p Printer
	switch f {
	case OutFormatText:
		p = NewText(util.GetEnvBoolDefault("DISABLE_COLORS", false))
	case OutFormatSimple:
		p = NewSimple()
	case OutFormatGitHubActions:
		p = NewGitHubActions()
	case OutFormatJSON:
		p = NewJSON()
	default:
		return p, fmt.Errorf("%s is not a valid printer type", f)
	}
	log.Debug().Str("printer", f).Msg("created new printer")
	return p, nil
}
