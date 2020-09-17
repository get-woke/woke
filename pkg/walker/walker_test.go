package walker

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalker_Walk(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	assert.DirExists(t, dir)

	for i := 0; i < 10; i++ {
		var newDir string
		if i%2 == 0 {
			newDir = filepath.Join(dir, strconv.Itoa(i))
		} else {
			newDir = filepath.Join(dir, ".git")
		}
		assert.NoError(t, os.MkdirAll(newDir, 0777))

		filename := filepath.Join(newDir, ".foo")
		file, err := os.Create(filename)
		assert.NoError(t, err)
		assert.NoError(t, file.Close())
	}

	err = Walk(dir, func(p string, typ os.FileMode) error {
		assert.False(t, isDotGit(p), "path should not be returned in walk: %s", p)
		return nil
	})
	assert.NoError(t, err)

}

func TestInSlice(t *testing.T) {
	tests := []struct {
		path      string
		assertion assert.BoolAssertionFunc
	}{
		{".git", assert.True},
		{".github", assert.False},
		{"/foo/bar/.git", assert.True},
		{"/foo/.git/bar", assert.False},
		{"/foo/.github", assert.False},
		{"foo/.git", assert.True},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			tt.assertion(t, isDotGit(tt.path))
		})
	}
}
