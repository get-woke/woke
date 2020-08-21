package util

import (
	"path/filepath"
)

func GetFilesInGlobs(globs []string) ([]string, error) {
	var files []string
	for _, glob := range globs {
		filesInGlob, err := filepath.Glob(glob)
		if err != nil {
			return files, err
		}
		files = append(files, filesInGlob...)
	}
	return files, nil
}
