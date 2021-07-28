package ignore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/get-woke/woke/pkg/util"
)

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
		domain = strings.Split(util.ParseRelativePathFromAbsolutePath(cwd, rootDir), string(os.PathSeparator))
		return rootDir, domain, nil
	}
	return cwd, domain, nil
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
