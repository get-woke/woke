package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/get-woke/woke/pkg/result"

	env "github.com/caitlinelfring/go-env-default"
	"github.com/rs/zerolog/log"
)

// Printer is an interface for printing FileResults
type Printer interface {
	Print(*result.FileResults) error
	Start()
	End()
	PrintSuccessExitMessage() bool
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

	// OutFormatSonarQube is an output format supported by SonarQube
	// https://docs.sonarqube.org/latest/analysis/generic-issue/
	OutFormatSonarQube = "sonarqube"

	// OutFormatCheckstyle outputs in checkstyle format.
	// https://github.com/checkstyle/checkstyle
	OutFormatCheckstyle = "checkstyle"
)

// OutFormats are all the available output formats. The first one should be the default
var OutFormats = []string{
	OutFormatText,
	OutFormatSimple,
	OutFormatGitHubActions,
	OutFormatJSON,
	OutFormatCheckstyle,
	OutFormatSonarQube,
}

// OutFormatsString is all OutFormats, as a comma-separated string
var OutFormatsString = strings.Join(OutFormats, ",")

// NewPrinter returns a valid new Printer from a string, or an error if the printer is invalid
func NewPrinter(f string, w io.Writer) (Printer, error) {
	var p Printer
	switch f {
	case OutFormatText:
		p = NewText(w, env.GetBoolDefault("DISABLE_COLORS", false))
	case OutFormatSimple:
		p = NewSimple(w)
	case OutFormatGitHubActions:
		p = NewGitHubActions(w)
	case OutFormatJSON:
		p = NewJSON(w)
	case OutFormatSonarQube:
		p = NewSonarQube(w)
	case OutFormatCheckstyle:
		p = NewCheckstyle(w)
	default:
		return p, fmt.Errorf("%s is not a valid printer type", f)
	}
	log.Debug().Str("printer", f).Msg("created new printer")
	return p, nil
}
