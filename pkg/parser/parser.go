package parser

import (
	"bufio"
	"os"
	"time"

	"github.com/caitlinelfring/woke/pkg/ignore"
	"github.com/caitlinelfring/woke/pkg/result"
	"github.com/caitlinelfring/woke/pkg/rule"
	"github.com/caitlinelfring/woke/pkg/util"
	"github.com/rs/zerolog/log"
)

// Parser parses files and finds lines that break rules
type Parser struct {
	Rules []*rule.Rule
}

// ParseFiles parses all files provided and returns the results
func (p *Parser) ParseFiles(files []string, ignorer *ignore.Ignore) (results []*result.FileResults) {
	parsable := WalkDirsWithIgnores(files, ignorer)

	for _, f := range parsable {
		if fileResult, err := p.ParseFile(f); fileResult != nil && err == nil {
			results = append(results, fileResult)
		}
	}

	return
}

// Parse reads the file and returns results of places where rules are broken
// this function will not close the file, that should be handled by the caller
func (p *Parser) Parse(file *os.File) (*result.FileResults, error) {
	start := time.Now()
	defer func() {
		log.Debug().
			Str("file", file.Name()).
			Dur("durationMS", time.Now().Sub(start)).
			Msg("finished Parse")
	}()

	results := &result.FileResults{
		Filename: file.Name(),
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		text := scanner.Text()
		for _, r := range p.Rules {
			lineResults := result.FindResults(r, text, line)
			results.Results = append(results.Results, lineResults...)
			// results.Push(lineResults...)
		}
		line++
	}

	return results, scanner.Err()
}

// ParseFile parses the files provided and returns the results
func (p *Parser) ParseFile(f string) (*result.FileResults, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err = util.IsTextFile(file); err != nil {
		return nil, err
	}

	r, err := p.Parse(file)
	if err != nil {
		return nil, err
	}

	if len(r.Results) == 0 {
		return nil, nil
	}

	return &result.FileResults{
		Filename: file.Name(),
		Results:  r.Results,
	}, nil
}
