package ignore

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	gitignore "github.com/sabhiram/go-gitignore"
)

// DefaultIgnores is the default list of file globs that will be ignored
var DefaultIgnores = []string{
	".git",
}

type Ignore struct {
	compiled *gitignore.GitIgnore
}

// NewIgnore produces an Ignore object, with compiled lines from .gitignore and DefaultIgnores
// which you can match files against
func NewIgnore(lines []string) (*Ignore, error) {
	compiled, err := compileIgnoreLines(lines)
	if err != nil {
		return nil, err
	}
	return &Ignore{compiled: compiled}, nil
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string) bool {
	return i.compiled.MatchesPath(f)
}

func compileIgnoreLines(lines []string) (*gitignore.GitIgnore, error) {
	lines = append(lines, DefaultIgnores...)
	lines = append(lines, readIgnoreFile(".gitignore")...)
	lines = append(lines, readIgnoreFile(".wokeignore")...)

	return gitignore.CompileIgnoreLines(lines...)
}

func readIgnoreFile(file string) []string {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		_event := log.Warn()
		if errors.Is(err, os.ErrNotExist) {
			_event = log.Debug()
		}
		_event.Err(err).Str("file", file).Msg("skipping ignorefile")
	}

	return strings.Split(string(buffer), "\n")
}
