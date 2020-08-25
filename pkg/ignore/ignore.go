package ignore

import (
	"io/ioutil"
	"strings"

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

	if buffer, err := ioutil.ReadFile(".gitignore"); err == nil {
		lines = append(lines, strings.Split(string(buffer), "\n")...)
	}

	return gitignore.CompileIgnoreLines(lines...)
}
