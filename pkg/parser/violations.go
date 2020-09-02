package parser

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/rs/zerolog/log"
)

func generateFileViolationsFromFilename(filename string, rules []*rule.Rule) (*result.FileResults, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return generateFileViolations(file, rules)
}

// generateFileViolations reads the file and returns results of places where rules are broken
// this function will not close the file, that should be handled by the caller
func generateFileViolations(file *os.File, rules []*rule.Rule) (*result.FileResults, error) {
	start := time.Now()
	defer func() {
		log.Debug().
			Str("file", file.Name()).
			Dur("durationMS", time.Since(start)).
			Msg("finished generateFileViolations")
	}()

	results := &result.FileResults{
		Filename: file.Name(),
	}

	reader := bufio.NewReader(file)

	line := 1
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, r := range rules {
			lineResults := result.FindResults(r, results.Filename, text, line)
			results.Results = append(results.Results, lineResults...)
		}

		line++
	}

	return results, nil
}
