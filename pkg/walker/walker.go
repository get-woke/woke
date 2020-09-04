package walker

import (
	"os"
	"path/filepath"

	"github.com/get-woke/fastwalk"
)

// Walk is a helper function that will automatically skip the `.git` directory.
// fastwalk is a fork of code that is a better, faster version of filepath.Walk.
// tl;dr since filepath.Walk get a complete FileInfo for every file,
// it's inherently slow. See https://github.com/golang/go/issues/16399
func Walk(root string, walkFn func(path string, typ os.FileMode) error) error {
	return fastwalk.Walk(root, func(path string, typ os.FileMode) error {
		path = filepath.Clean(path)

		if typ.IsDir() && isDotGit(path) {
			return filepath.SkipDir
		}

		return walkFn(path, typ)
	})
}

func isDotGit(path string) bool {
	return filepath.Base(path) == ".git"
}
