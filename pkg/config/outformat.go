package config

import (
	"fmt"
	"strings"

	"github.com/get-woke/woke/pkg/printer"
	"github.com/get-woke/woke/pkg/util"
	"github.com/rs/zerolog/log"
)

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

func CreatePrinter(f string) (printer.Printer, error) {
	var p printer.Printer
	switch f {
	case OutFormatText:
		p = printer.NewText(util.GetEnvBoolDefault("DISABLE_COLORS", false))
	case OutFormatSimple:
		p = printer.NewSimple()
	case OutFormatGitHubActions:
		p = printer.NewGitHubActions()
	default:
		return p, fmt.Errorf("%s is not a valid printer type", f)
	}
	log.Debug().Str("printer", f).Msg("created new printer")
	return p, nil
}
