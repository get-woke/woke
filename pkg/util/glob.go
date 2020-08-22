package util

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// GetFilesInGlobs returns all known files in the provided globs using
// filepath.Glob and filepath.Walk
func GetFilesInGlobs(globs []string) ([]string, bool, error) {
	var files []string
	useAbsolutePath := false
	for _, glob := range globs {
		if filepath.IsAbs(glob) {
			useAbsolutePath = true
		}
	}
	for _, glob := range globs {
		filesInGlob, err := filepath.Glob(glob)
		if err != nil {
			return files, useAbsolutePath, err
		}

		for _, p := range filesInGlob {
			err := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
				if useAbsolutePath {
					abs, err := filepath.Abs(path)
					if err != nil {
						return err
					}
					files = append(files, abs)
				} else {
					files = append(files, path)
				}
				return err
			})

			if err != nil {
				log.Error().
					Err(err).
					Msgf("error walking %s", p)
			}
		}
	}
	return files, useAbsolutePath, nil
}
