package util

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func GetFilesInGlobs(globs []string) ([]string, error) {
	var files []string
	for _, glob := range globs {
		filesInGlob, err := filepath.Glob(glob)
		if err != nil {
			return files, err
		}

		for _, p := range filesInGlob {
			err := filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
				files = append(files, path)
				return err
			})

			if err != nil {
				log.Error().
					Err(err).
					Msgf("error walking %s", p)
			}
		}
	}
	return files, nil
}
