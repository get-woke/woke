package parser

import (
	"bufio"
	"os"
	"time"

	"github.com/caitlinelfring/woke/pkg/ignore"
	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/caitlinelfring/woke/pkg/util"
	"github.com/rs/zerolog/log"
)

// Parser parses files and finds lines that break rules
type Parser struct {
	Rules []*rule.Rule
}

// ParseFiles parses all files provided and returns the results
func (p *Parser) ParseFiles(files []string, ignorer *ignore.Ignore) rule.Results {
	parsable := ParsableFiles(files, ignorer)
	return p.parseFiles(parsable.Files)
}

// Parse reads the file and returns results of places where rules are broken
// this function will not close the file, that should be handled by the caller
func (p *Parser) Parse(file *os.File) (results rule.Results, err error) {
	start := time.Now()
	defer func() {
		log.Debug().
			Str("file", file.Name()).
			Dur("durationMS", time.Now().Sub(start)).
			Msg("finished Parse")
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		text := scanner.Text()
		for _, r := range p.Rules {
			lineResults := r.FindResults(text, file.Name(), line)
			results.Push(lineResults...)
		}

		line++
	}

	return results, scanner.Err()
}

func (p *Parser) parseFiles(files []string) rule.Results {
	results := rule.Results{}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			log.Error().Err(err).Str("file", f).Msg("could not open file")
			continue
		}
		defer file.Close()

		if err = util.IsTextFile(file); err != nil {
			log.Debug().Err(err).Str("file", file.Name()).Msg("not a text file")
			continue
		}

		r, err := p.Parse(file)
		if err != nil {
			log.Debug().Err(err).Msg("parser failed")
			continue
		}

		results.Push(r.Results...)
	}

	return results
}
