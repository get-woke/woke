package parser

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/get-woke/woke/pkg/result"

	"github.com/rs/zerolog/log"
)

func (p *Parser) generateFileViolationsFromFilename(filename string) (*result.FileResults, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return p.generateFileViolations(file)
}

// generateFileViolations reads the file and returns results of places where rules are broken
// this function will not close the file, that should be handled by the caller
func (p *Parser) generateFileViolations(file *os.File) (*result.FileResults, error) {
	filename := filepath.ToSlash(file.Name())
	start := time.Now()
	defer func() {
		log.Debug().
			TimeDiff("durationMS", time.Now(), start).
			Str("file", filename).
			Msg("finished processing violations")
	}()

	results := &result.FileResults{
		Filename: filename,
	}

	// Check for violations in the filename itself
	for _, pathResult := range result.MatchPathRules(p.Rules, file.Name()) {
		results.Results = append(results.Results, pathResult)
	}

	reader := bufio.NewReader(file)

	line := 1

Loop:
	for {
		switch text, err := reader.ReadString('\n'); {
		case err == nil || (err == io.EOF && text != ""):
			text = strings.TrimSuffix(text, "\n")

			for _, r := range p.Rules {
				lineResults := result.FindResults(r, results.Filename, text, line)
				results.Results = append(results.Results, lineResults...)
			}

			line++
		case err == io.EOF:
			break Loop
		case err != nil:
			return nil, err
		}
	}

	return results, nil
}
