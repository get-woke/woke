package parser

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	filename := filepath.ToSlash(file.Name())
	start := time.Now()
	defer func() {
		log.Debug().
			TimeDiff("durationMS", time.Now(), start).
			Str("file", filename).
			Msg("finished generateFileViolations")
	}()

	results := &result.FileResults{
		Filename: filename,
	}

	// Check for violations in the filename itself
	for _, pathResult := range result.MatchPathRules(rules, file.Name()) {
		results.Results = append(results.Results, pathResult)
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

		// text will have a trailing new line, trim it
		text = strings.TrimSuffix(text, "\n")

		for _, r := range rules {
			lineResults := result.FindResults(r, results.Filename, text, line)
			results.Results = append(results.Results, lineResults...)
		}

		line++
	}

	return results, nil
}
