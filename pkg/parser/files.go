package parser

import (
	"os"
	"path/filepath"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/rs/zerolog/log"
)

const defaultPath = "."

// WalkDirsWithIgnores returns a list of files that can be parsed after the ignorer has
// excluded files that should be ignored
func WalkDirsWithIgnores(paths []string, ignorer *ignore.Ignore) []string {
	if len(paths) == 0 {
		paths = []string{defaultPath}
	}

	files, _ := WalkDirs(paths, ignorer)

	return files
}

// WalkDirs returns all known files in the provided paths using filepath.Walk
func WalkDirs(paths []string, ignorer *ignore.Ignore) ([]string, error) {
	var files []string

	for _, p := range paths {
		err := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignore directories and files that match the ignorer
			if !f.IsDir() && !ignorer.Match(path) {
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
