package ignore

import (
	"os"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
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
	if strings.LastIndex(workingDir, string(os.PathSeparator)) != len(workingDir)-1 {
		workingDir += string(os.PathSeparator)
	}

	res := strings.SplitN(absPath, workingDir, 2)
	if len(res) > 1 {
		x := util.FilterEmptyStrings(strings.Split(res[1], string(os.PathSeparator)))
		return x
	}
	return []string{}
}

func GetRootGitDir(path string) (filesystem billy.Filesystem, err error) {
	opt := &git.PlainOpenOptions{DetectDotGit: true}
	var repo *git.Repository
	repo, err = git.PlainOpenWithOptions(path, opt)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			log.Debug().Str("reason", err.Error()).Msg("Could Not Find Root Git Folder")
			filesystem = osfs.New(path)
		} else {
			return
		}
	}
	if repo != nil {
		var worktree *git.Worktree
		worktree, err = repo.Worktree()
		if err != nil {
			return
		}
		filesystem = worktree.Filesystem
	}
	return
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

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	ps, err := readPatterns(filesystem, nil)
	if err != nil {
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
