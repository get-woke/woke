package ignore

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/get-woke/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/rs/zerolog/log"

	"github.com/get-woke/woke/pkg/util"
)

// Ignore is a gitignore-style object to ignore files/directories
type Ignore struct {
	matcher gitignore.Matcher
}

var defaultIgnoreFiles = []string{
	".gitignore",
	".ignore",
	".wokeignore",
}

// Given a workingDir (example: /root/proj/subDir/curDir)
// and given a gitRoot (example: /root/proj)
// it will return []string{"subDir", "curDir"}.
// This is the domain (path from git root) that ignore rules will apply to.
func getDomainFromWorkingDir(workingDir, gitRoot string) []string {
	// if working directory does not end with a slash, add it
	if !strings.HasSuffix(gitRoot, string(os.PathSeparator)) {
		gitRoot += string(os.PathSeparator)
	}

	res := strings.SplitN(workingDir, gitRoot, 2)
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
func NewIgnore(filesystem billy.Filesystem, lines []string, disableNestedIgnores bool) (ignore *Ignore, err error) {
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

	// if opted-out of nested wokeignore traversal, only use top-level ignore files
	if disableNestedIgnores {
		for _, filename := range defaultIgnoreFiles {
			lines = append(lines, readIgnoreFile(filesystem, filename)...)
		}
	} else {
		if ps, err = gitignore.ReadPatterns(filesystem, nil, defaultIgnoreFiles); err != nil {
			return
		}
	}

	// get domain for git ignore rules supplied from the lines argument
	domain := getDomainFromWorkingDir(cwd, filesystem.Root())
	for _, line := range lines {
		pattern := gitignore.ParsePattern(line, domain)
		ps = append(ps, pattern)
	}

	return &Ignore{
		matcher: gitignore.NewMatcher(ps),
	}, nil
}

func readIgnoreFile(filesystem billy.Filesystem, file string) []string {
	openFile, err := filesystem.Open(file)

	if err != nil {
		_event := log.Warn()
		if errors.Is(err, os.ErrNotExist) {
			_event = log.Debug()
		}
		_event.Err(err).Str("file", file).Msg("skipping ignorefile")
		return []string{}
	}

	defer openFile.Close()
	var buffer []byte
	buffer, err = io.ReadAll(bufio.NewReader(openFile))

	if err != nil {
		_event := log.Warn()
		_event.Err(err).Str("file", file).Msg("skipping ignorefile")
		return []string{}
	}

	log.Debug().Str("file", file).Msg("adding ignorefile")
	return strings.Split(strings.TrimSpace(string(buffer)), "\n")
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string, isDir bool) bool {
	parts := util.FilterEmptyStrings(strings.Split(f, string(os.PathSeparator)))
	return i.matcher.Match(parts, isDir)
}
