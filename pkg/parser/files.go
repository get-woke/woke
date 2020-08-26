package parser

import (
	"os"
	"path/filepath"

	"github.com/caitlinelfring/woke/pkg/ignore"
	"github.com/rs/zerolog/log"
)

const defaultPath = "."

// WalkDirsWithIgnores returns a list of files that can be parsed after the ignorer has
// excluded files that should be ignored
func WalkDirsWithIgnores(paths []string, ignorer *ignore.Ignore) (files []string) {
	if len(paths) == 0 {
		paths = []string{defaultPath}
	}

	allFiles, _ := WalkDirs(paths)
	if ignorer == nil {
		return allFiles
	}

	for _, f := range allFiles {
		if ignorer.Match(f) {
			continue
		}
		files = append(files, f)
	}

	return
}

// WalkDirs returns all known files in the provided paths using filepath.Walk
func WalkDirs(paths []string) ([]string, error) {
	var files []string

	for _, p := range paths {
		err := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
			// Ignore directories
			if !f.IsDir() {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			log.Error().
				Err(err).
				Msgf("error walking %s", p)
		}
	}

	return files, nil
}
