package ignore

// These functions are loosely modified from https://github.com/go-git/go-git/blob/master/plumbing/format/gitignore/dir.go
// in order to support .wokeignore files as well as .gitignore files

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/rs/zerolog/log"
)

const (
	commentPrefix = "#"
	gitDir        = ".git"
)

var defaultIgnoreFiles = []string{
	".gitignore",
	".ignore",
	".wokeignore",
	".git/info/exclude",
}

// readIgnoreFile reads a specific git ignore file.
func readIgnoreFile(fs billy.Filesystem, path []string, ignoreFile string) (ps []gitignore.Pattern, err error) {
	ps = []gitignore.Pattern{}
	var f billy.File
	ignoreFilePath := fs.Join(append(path, ignoreFile)...)
	f, err = fs.Open(ignoreFilePath)
	if err != nil {
		_event := log.Warn()
		if errors.Is(err, os.ErrNotExist) {
			_event = log.Debug()
			err = nil
		}
		_event.Err(err).Str("file", ignoreFilePath).Msg("skipping ignorefile")
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if !strings.HasPrefix(s, commentPrefix) && len(strings.TrimSpace(s)) > 0 {
			ps = append(ps, gitignore.ParsePattern(s, path))
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	return
}

// readPatterns reads gitignore patterns recursively traversing through the directory
// structure. The result is in the ascending order of priority (last higher).
func readPatterns(fs billy.Filesystem, path []string) (ps []gitignore.Pattern, err error) {
	ps = []gitignore.Pattern{}
	for _, filename := range defaultIgnoreFiles {
		var subps []gitignore.Pattern
		if subps, err = readIgnoreFile(fs, path, filename); err != nil {
			return ps, err
		}
		if len(subps) > 0 {
			ps = append(ps, subps...)
		}
	}

	var fis []os.FileInfo
	if fis, err = fs.ReadDir(fs.Join(path...)); err != nil {
		return ps, err
	}

	for _, fi := range fis {
		if fi.IsDir() && fi.Name() != gitDir {
			var subps []gitignore.Pattern
			if subps, err = readPatterns(fs, append(path, fi.Name())); err != nil {
				return ps, err
			}

			if len(subps) > 0 {
				ps = append(ps, subps...)
			}
		}
	}

	return ps, nil
}
