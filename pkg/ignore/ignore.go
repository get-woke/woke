package ignore

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/rs/zerolog/log"

	"github.com/get-woke/woke/pkg/util"
)

// Ignore is a gitignore-style object to ignore files/directories
type Ignore struct {
	matcher gitignore.Matcher
}

const (
	commentPrefix = "#"
	gitDir        = ".git"
)

var defaultIgnoreFiles = []string{
	".gitignore",
	".wokeignore",
}

// readIgnoreFile reads a specific git ignore file.
func readIgnoreFile(fs billy.Filesystem, path []string, ignoreFile string) (ps []gitignore.Pattern, err error) {
	f, err := fs.Open(fs.Join(append(path, ignoreFile)...))
	if err != nil {
		_event := log.Warn()
		if errors.Is(err, os.ErrNotExist) {
			_event = log.Debug()
			err = nil
		}
		_event.Err(err).Str("file", fs.Join(append(path, ignoreFile)...)).Msg("skipping ignorefile")
		return []gitignore.Pattern{}, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, commentPrefix) && len(strings.TrimSpace(s)) > 0 {
			ps = append(ps, gitignore.ParsePattern(s, path))
		}
	}

	return
}

// readPatterns reads gitignore patterns recursively traversing through the directory
// structure. The result is in the ascending order of priority (last higher).
func readPatterns(fs billy.Filesystem, path []string) (ps []gitignore.Pattern, err error) {
	ps = []gitignore.Pattern{}
	for _, filename := range defaultIgnoreFiles {
		var subps []gitignore.Pattern
		subps, err = readIgnoreFile(fs, path, filename)
		if err != nil {
			return ps, err
		}
		if len(subps) > 0 {
			ps = append(ps, subps...)
		}
	}

	var fis []os.FileInfo
	fis, err = fs.ReadDir(fs.Join(path...))
	if err != nil {
		return ps, err
	}

	for _, fi := range fis {
		if fi.IsDir() && fi.Name() != gitDir {
			var subps []gitignore.Pattern
			subps, err = readPatterns(fs, append(path, fi.Name()))
			if err != nil {
				return ps, err
			}

			if len(subps) > 0 {
				ps = append(ps, subps...)
			}
		}
	}

	return ps, nil
}

// NewIgnore produces an Ignore object, with compiled lines from defaultIgnoreFiles
// which you can match files against
func NewIgnore(lines []string, findRootDir bool) *Ignore {
	start := time.Now()
	defer func() {
		log.Debug().
			TimeDiff("durationMS", time.Now(), start).
			Msg("finished compiling ignores")
	}()

	rootDir, domain, err := getRootFileSystem(findRootDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	rootFs := osfs.New(rootDir)
	currentPath := []string{"."}
	ps, err := readPatterns(rootFs, currentPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, line := range lines {
		pattern := gitignore.ParsePattern(line, domain)
		ps = append(ps, pattern)
	}

	ignorer := Ignore{
		matcher: gitignore.NewMatcher(ps),
	}

	return &ignorer
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string, isDir bool) bool {
	parts := util.FilterEmptyStrings(strings.Split(f, string(os.PathSeparator)))
	return i.matcher.Match(parts, isDir)
}
