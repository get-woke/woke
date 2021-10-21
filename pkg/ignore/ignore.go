package ignore

import (
	"os"
	"path/filepath"
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

// given a absolute path (example: /runner/folder/root/data/effx.yaml)
// and given a workingDir (example: root)
// it will return data/effx.yaml
func getDomainFromWorkingDir(absPath, workingDir string) []string {
	// if working directory does not end with a slash, add it
	if !strings.HasSuffix(workingDir, string(os.PathSeparator)) {
		workingDir += string(os.PathSeparator)
	}

	res := strings.SplitN(absPath, workingDir, 2)
	if len(res) > 1 {
		x := util.FilterEmptyStrings(strings.Split(res[1], string(os.PathSeparator)))
		return x
	}
	return []string{}
}

func GetRootGitDir(workingDir string) (filesystem billy.Filesystem, err error) {
	dir, err := filepath.Abs(workingDir)
	if err != nil {
		return
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return osfs.New(dir), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			log.Debug().Msg("Could Not Find Root Git Folder")
			return osfs.New(workingDir), nil
		}
		dir = parent
	}
}

// NewIgnore produces an Ignore object, with compiled lines from defaultIgnoreFiles
// which you can match files against
func NewIgnore(filesystem billy.Filesystem, lines []string) (ignore *Ignore, err error) {
	start := time.Now()
	defer func() {
		log.Debug().
			TimeDiff("durationMS", time.Now(), start).
			Msg("finished compiling ignores")
	}()

	var cwd string
	if cwd, err = os.Getwd(); err != nil {
		return
	}

	var ps []gitignore.Pattern
	if ps, err = readPatterns(filesystem, nil); err != nil {
		return
	}

	domain := getDomainFromWorkingDir(cwd, filesystem.Root())
	for _, line := range lines {
		pattern := gitignore.ParsePattern(line, domain)
		ps = append(ps, pattern)
	}

	ignore = &Ignore{
		matcher: gitignore.NewMatcher(ps),
	}

	return
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string, isDir bool) bool {
	parts := util.FilterEmptyStrings(strings.Split(f, string(os.PathSeparator)))
	return i.matcher.Match(parts, isDir)
}
