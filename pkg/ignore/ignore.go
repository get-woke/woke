package ignore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/rs/zerolog/log"
)

// Ignore is a gitignore-style object to ignore files/directories
type Ignore struct {
	matcher Matcher
}

// given a absolute path (example: /runner/folder/root/data/effx.yaml)
// and given a workingDir (example: root)
// it will return data/effx.yaml
func parseRelativePathFromAbsolutePath(absoluteDir, workingDir string) string {
	// if working directory does not end with a slash, add it
	if strings.LastIndex(workingDir, "/") != len(workingDir)-1 {
		workingDir += "/"
	}

	res := strings.Split(absoluteDir, workingDir)
	if len(res) > 1 {
		return res[1]
	}
	return ""
}

func getRootFileSystem(findRootDir bool) (string, []string, error) {
	var domain []string
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return "", domain, err
	}
	if findRootDir {
		rootDir, err := detectGitPath(".")
		if err != nil {
			fmt.Println(err)
			return "", domain, err
		}
		domain = strings.Split(parseRelativePathFromAbsolutePath(cwd, rootDir), string(os.PathSeparator))
		return rootDir, domain, nil
	}
	return cwd, domain, nil
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
	ps, err := ReadPatterns(rootFs, currentPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, line := range lines {
		pattern := ParsePattern(line, domain)
		ps = append(ps, pattern)
	}

	ignorer := Ignore{
		matcher: NewMatcher(ps),
	}

	return &ignorer
}

func filterEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// Match returns true if the provided file matches any of the defined ignores
func (i *Ignore) Match(f string, isDir bool) bool {
	parts := filterEmptyStrings(strings.Split(f, string(os.PathSeparator)))
	return i.matcher.Match(parts, isDir)
}

func detectGitPath(path string) (string, error) {
	// normalize the path
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		fi, err := os.Stat(filepath.Join(path, gitDir))
		if err == nil {
			if !fi.IsDir() {
				return "", fmt.Errorf("%v exist but is not a directory", gitDir)
			}
			return path, nil
		}
		if !os.IsNotExist(err) {
			// unknown error
			return "", err
		}

		// detect bare repo
		ok, err := isGitDir(path)
		if err != nil {
			return "", err
		}
		if ok {
			return path, nil
		}

		if parent := filepath.Dir(path); parent == path {
			return "", fmt.Errorf(".git not found")
		} else {
			path = parent
		}
	}
}

func isGitDir(path string) (bool, error) {
	markers := []string{"HEAD", "objects", "refs"}

	for _, marker := range markers {
		_, err := os.Stat(filepath.Join(path, marker))
		if err == nil {
			continue
		}
		if !os.IsNotExist(err) {
			// unknown error
			return false, err
		} else {
			return false, nil
		}
	}

	return true, nil
}
