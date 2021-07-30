package ignore

import (
	"fmt"
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

	res := strings.Split(absPath, workingDir)
	if len(res) > 1 {
		return util.FilterEmptyStrings(strings.Split(res[1], string(os.PathSeparator)))
	}
	return []string{}
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

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// DetectDotGit will traverse parent directories until it finds the root git dir if true
	opt := &git.PlainOpenOptions{DetectDotGit: findRootDir}
	repo, err := git.PlainOpenWithOptions("", opt)
	var filesystem billy.Filesystem
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			log.Debug().Str("reason", err.Error()).Msg("Could Not Find Root Git Folder")
			filesystem = osfs.New(cwd)
		} else {
			fmt.Println(err)
			return nil
		}
	}
	if repo != nil {
		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		filesystem = worktree.Filesystem
	}

	ps, err := readPatterns(filesystem, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	domain := getDomainFromWorkingDir(cwd, filesystem.Root())

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
