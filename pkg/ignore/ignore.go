package ignore

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	gitignore "github.com/get-woke/go-gitignore"
	"github.com/rs/zerolog/log"
)

// DefaultIgnores is the default list of file globs that will be ignored
var DefaultIgnores = []string{
	".git",
}

// Ignore is a gitignore-style object to ignore files/directories
type Ignore struct {
	matcher *gitignore.GitIgnore
}

// NewIgnore produces an Ignore object, with compiled lines from .gitignore and DefaultIgnores
// which you can match files against
func NewIgnore(lines []string, pathsForGitIgnores []string) *Ignore {
	start := time.Now()
	defer func() {
		log.Debug().
			Dur("durationMS", time.Since(start)).
			Msg("finished compiling ignores")
	}()

	lines = append(lines, DefaultIgnores...)
	lines = append(lines, readIgnoreFile(".gitignore")...)
	lines = append(lines, readIgnoreFile(".wokeignore")...)

	ignorer := Ignore{
		matcher: gitignore.CompileIgnoreLines(lines...),
	}

	// FIXME: This is very costly with large directories with a lot of files, disabled for now
	// ignorer.AddIgnoreFiles(".gitignore", pathsForGitIgnores)
	// ignorer.AddIgnoreFiles(".wokeignore", pathsForGitIgnores)

	return &ignorer
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string) bool {
	return i.matcher.MatchesPath(f)
}

// AddIgnoreFiles walks each path provided in search of any files that match ignoreName
// and add the contents of those files to the gitignore matcher
// NOTE: this is very costly in large directories and should be used with caution
func (i *Ignore) AddIgnoreFiles(ignoreName string, paths []string) {
	lines := addRecursiveGitIgnores(ignoreName, paths)
	i.matcher.AddPatternsFromLines(lines...)
}

// addRecursiveGitIgnores uses filepath.Walk to walk each path, search for a file that matches
// ignoreName and reads each file's lines
// NOTE: this is very costly in large directories and should be used with caution
func addRecursiveGitIgnores(ignoreName string, paths []string) (lines []string) {
	for _, path := range paths {
		_ = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.Mode().IsRegular() && info.Name() == ignoreName {
				newLines := append(readIgnoreFile(p), p)
				lines = append(lines, newLines...)
			}

			return nil
		})
	}

	return
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
	rawLines := strings.Split(strings.TrimSpace(string(buffer)), "\n")

	// Pre-allocate the slice
	lines := make([]string, len(rawLines))
	// Here, we are prefixing each line with the base of the ignore file
	// to suppose ignore files in subdirectories
	for i := range rawLines {
		lines[i] = filepath.Join(filepath.Dir(file), rawLines[i])
	}

	return lines
}
