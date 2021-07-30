package ignore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestIgnore_Match(t *testing.T) {
	i := NewIgnore([]string{"my/files/*"}, false)
	assert.NotNil(t, i)

	// Test if rules with backslashes match on windows
	assert.False(t, i.Match("not/foo", false))
	assert.True(t, i.Match("my/files/file1", false))
	assert.False(t, i.Match("my/files", false))

	assert.False(t, i.Match(filepath.Join("not", "foo"), false))
	assert.True(t, i.Match(filepath.Join("my", "files", "file1"), false))
	assert.False(t, i.Match(filepath.Join("my", "files"), false))
}

// Test all default ignore files, except for .git/info/exclude, since
// that uses a .git directory that we cannot check in.
func TestIgnoreDefaultIgoreFiles_Match(t *testing.T) {
	// Temporarily change into testdata directojry for this test
	// since paths are relative
	err := os.Chdir("testdata")
	assert.NoError(t, err)
	t.Cleanup(func() {
		err = os.Chdir("..")
		assert.NoError(t, err)
	})

	i := NewIgnore([]string{"*.FROMARGUMENT"}, false)
	assert.NotNil(t, i)

	assert.False(t, i.Match("notfoo", false))
	assert.True(t, i.Match("test.FROMARGUMENT", false)) // From .gitignore
	assert.True(t, i.Match("test.DS_Store", false))     // From .gitignore
	assert.True(t, i.Match("test.IGNORE", false))       // From .ignore
	assert.True(t, i.Match("test.WOKEIGNORE", false))   // From .wokeignore
	assert.False(t, i.Match("test.NOTIGNORED", false))  // From .notincluded - making sure only default are included
}

func TestReadIgnoreFile(t *testing.T) {
	rootFs := osfs.New(".")
	ignoreFilePath := []string{"testdata"}
	ignoreLines, _ := readIgnoreFile(rootFs, ignoreFilePath, ".gitignore")
	patterns := []gitignore.Pattern{gitignore.ParsePattern("*.DS_Store", []string{"testdata"})}
	assert.Equal(t, patterns, ignoreLines)

	noIgnoreLines, _ := readIgnoreFile(rootFs, []string{}, ".gitignore")
	assert.Equal(t, []gitignore.Pattern{}, noIgnoreLines)
}
