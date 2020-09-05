package ignore

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gitignore "github.com/get-woke/go-gitignore"
	"github.com/get-woke/woke/pkg/util"
	"github.com/get-woke/woke/pkg/walker"
	"github.com/rs/zerolog/log"
)

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
			TimeDiff("durationMS", time.Now(), start).
			Msg("finished compiling ignores")
	}()

	lines = append(lines, readIgnoreFile(".gitignore")...)
	lines = append(lines, readIgnoreFile(".wokeignore")...)

	ignorer := Ignore{
		matcher: gitignore.CompileIgnoreLines(lines...),
	}

	// FIXME: This is very costly with large directories with a lot of files, disabled for now
	// ignorer.AddIgnoreFiles(pathsForGitIgnores, ".gitignore", ".wokeignore")

	return &ignorer
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string) bool {
	return i.matcher.MatchesPath(f)
}

// AddIgnoreFiles walks each path provided in search of any files that match ignoreName
// and add the contents of those files to the gitignore matcher
// NOTE: this is very costly in large directories and should be used with caution
func (i *Ignore) AddIgnoreFiles(paths []string, ignoreNames ...string) {
	lines := addRecursiveGitIgnores(ignoreNames, paths)
	i.matcher.AddPatternsFromLines(lines...)
}

// addRecursiveGitIgnores walks each path, search for a file that matches
// ignoreName and reads each file's lines
// NOTE: this is very costly in large directories and should be used with caution
func addRecursiveGitIgnores(ignoreNames []string, paths []string) (lines []string) {
	start := time.Now()
	defer func() {
		log.Debug().
			Strs("files", ignoreNames).
			TimeDiff("durationMS", time.Now(), start).
			Msg("finished walk for ignore files")
	}()

	var wg sync.WaitGroup
	wg.Add(len(paths))

	ch := make(chan []string)

	for _, path := range paths {
		go func(path string, ignoreNames []string) {
			defer wg.Done()

			_ = walker.Walk(path, func(p string, info os.FileMode) error {
				if info.Perm().IsRegular() && util.InSlice(filepath.Base(p), ignoreNames) {
					newLines := append(readIgnoreFile(p), p)
					ch <- newLines
				}

				return nil
			})
		}(path, ignoreNames)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for c := range ch {
		lines = append(lines, c...)
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
