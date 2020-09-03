package ignore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestIgnoreMatch(t *testing.T) {
	i := NewIgnore([]string{"my/files/*"}, []string{})
	assert.NotNil(t, i)

	assert.False(t, i.Match("not/foo"))
	assert.True(t, i.Match("my/files/file1"))
	assert.False(t, i.Match("my/files"))
}

func TestIgnore_AddIgnoreFiles(t *testing.T) {
	i := NewIgnore([]string{"my/files/*"}, []string{"."})
	i.AddIgnoreFiles(".gitignore", []string{"testdata"})

	assert.True(t, i.Match("testdata/.gitignore"))
	assert.True(t, i.Match("testdata/.DS_Store"))
	assert.False(t, i.Match(".DS_Store"))
	assert.True(t, i.Match("my/files/file.txt"))
	assert.False(t, i.Match("my/files"))
}

func TestAddRecursiveGitIgnores(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(t, err)
	assert.DirExists(t, dir)
	expected := make([]string, 0)
	lastDir := dir
	for i := 0; i < 5; i++ {
		lastDir = filepath.Join(lastDir, strconv.Itoa(i))
		assert.NoError(t, os.MkdirAll(lastDir, 0777))

		ignoreFilename := filepath.Join(lastDir, ".gitignore")
		file, err := os.Create(ignoreFilename)
		assert.NoError(t, err)

		content := fmt.Sprintf("%d.txt", i)
		_, _ = file.WriteString(content)
		assert.NoError(t, file.Close())

		expected = append(expected, filepath.Join(lastDir, content))
		expected = append(expected, ignoreFilename)
	}
	defer os.RemoveAll(dir)
	lines := addRecursiveGitIgnores(".gitignore", []string{dir})

	assert.EqualValues(t, expected, lines)
}

func BenchmarkIgnoreAddIgnoreFiles(b *testing.B) {
	dir, err := ioutil.TempDir(os.TempDir(), "")
	assert.NoError(b, err)
	assert.DirExists(b, dir)
	expected := []string{}

	for i := 0; i < b.N; i++ {
		newDir := filepath.Join(dir, strconv.Itoa(i))
		assert.NoError(b, os.MkdirAll(newDir, 0777))

		ignoreFilename := filepath.Join(newDir, ".gitignore")
		file, err := os.Create(ignoreFilename)
		assert.NoError(b, err)

		content := fmt.Sprintf("%d.txt", i)
		_, _ = file.WriteString(content)
		assert.NoError(b, file.Close())

		expected = append(expected, ignoreFilename)
		expected = append(expected, filepath.Join(newDir, content))
	}

	defer os.RemoveAll(dir)
	lines := addRecursiveGitIgnores(".gitignore", []string{dir})
	sort.Strings(lines)
	sort.Strings(expected)

	assert.EqualValues(b, expected, lines)
}
