package printer

import (
	"fmt"
	"strings"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/util"
	"github.com/rs/zerolog/log"
)

// Printer is an interface for printing FileResults
type Printer interface {
	Print(*result.FileResults) error
}

const (
	OutFormatText          = "text"
	OutFormatSimple        = "simple"
	OutFormatGitHubActions = "github-actions"
)

var OutFormats = []string{
	OutFormatText,
	OutFormatSimple,
	OutFormatGitHubActions,
}

var OutFormatsString = strings.Join(OutFormats, ",")

func NewPrinter(f string) (Printer, error) {
	var p Printer
	switch f {
	case OutFormatText:
		p = NewText(util.GetEnvBoolDefault("DISABLE_COLORS", false))
	case OutFormatSimple:
		p = NewSimple()
	case OutFormatGitHubActions:
		p = NewGitHubActions()
	default:
		return p, fmt.Errorf("%s is not a valid printer type", f)
	}
	log.Debug().Str("printer", f).Msg("created new printer")
	return p, nil
}
