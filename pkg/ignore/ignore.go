package ignore

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	gitignore "github.com/get-woke/go-gitignore"
	"github.com/rs/zerolog/log"
)

// Ignore is a gitignore-style object to ignore files/directories
type Ignore struct {
	matcher *gitignore.GitIgnore
}

// NewIgnore produces an Ignore object, with compiled lines from .gitignore and DefaultIgnores
// which you can match files against
func NewIgnore(lines []string) *Ignore {
	start := time.Now()
	defer func() {
		log.Debug().
			TimeDiff("durationMS", time.Now(), start).
			Msg("finished compiling ignores")
	}()

	lines = append(lines, readIgnoreFile(".gitignore")...)
	lines = append(lines, readIgnoreFile(".wokeignore")...)

	ignorer := Ignore{
		matcher: gitignore.CompileIgnoreLines(lines...),
	}

	return &ignorer
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string) bool {
	return i.matcher.MatchesPath(f)
}

func readIgnoreFile(file string) []string {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		_event := log.Warn()
		if errors.Is(err, os.ErrNotExist) {
			_event = log.Debug()
		}
		_event.Err(err).Str("file", file).Msg("skipping ignorefile")
		return []string{}
	}

	log.Debug().Str("file", file).Msg("adding ignorefile")

	return strings.Split(strings.TrimSpace(string(buffer)), "\n")
}
