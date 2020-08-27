package config

import (
	"strings"

	"github.com/caitlinelfring/woke/pkg/printer"
	"github.com/caitlinelfring/woke/pkg/util"
	"github.com/rs/zerolog/log"
)

const (
	OutFormatText   = "text"
	OutFormatSimple = "simple"
)

var OutFormats = []string{
	OutFormatText,
	OutFormatSimple,
}

var OutFormatsString = strings.Join(OutFormats, ",")

func CreatePrinter(f string) printer.Printer {
	var p printer.Printer
	switch f {
	case OutFormatText:
		p = printer.NewText(util.GetEnvBoolDefault("ENABLE_COLORS", true))
	case OutFormatSimple:
		p = printer.NewSimple()
	}
	log.Debug().Str("printer", f).Msg("created new printer")
	return p
}
