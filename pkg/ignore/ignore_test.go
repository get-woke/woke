package ignore

import (
	"os"
	"runtime"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestIgnore_Match(t *testing.T) {
	i := NewIgnore([]string{"my/files/*"})
	assert.NotNil(t, i)

	assert.False(t, i.Match("not/foo"))
	assert.True(t, i.Match("my/files/file1"))
	assert.False(t, i.Match("my/files"))

	if runtime.GOOS == "windows" {
		assert.False(t, i.Match(`not\foo`))
		assert.True(t, i.Match(`my\files\file1`))
		assert.False(t, i.Match(`my\files`))
	}
}

// Test all default ignore files, except for .git/info/exclude, since
// that uses a .git directory that we cannot check in.
func TestIgnoreDefaultIgoreFiles_Match(t *testing.T) {

	// Change into testdata directory for this test
	// since paths are relative
	err := os.Chdir("testdata")
	if err != nil {
		panic(err)
	}

	i := NewIgnore([]string{"*.FROMARGUMENT"})
	assert.NotNil(t, i)

	assert.False(t, i.Match("notfoo"))
	assert.True(t, i.Match("test.FROMARGUMENT")) // From .gitignore
	assert.True(t, i.Match("test.DS_Store"))     // From .gitignore
	assert.True(t, i.Match("test.IGNORE"))       // From .ignore
	assert.True(t, i.Match("test.WOKEIGNORE"))   // From .wokeignore
	assert.False(t, i.Match("test.NOTIGNORED"))  // From .notincluded - making sure only default are included

	// Change directory back to original directory
	err = os.Chdir("..")
	if err != nil {
		panic(err)
	}
}

func TestReadIgnoreFile(t *testing.T) {
	ignoreLines := readIgnoreFile("testdata/.gitignore")
	assert.Equal(t, []string{"*.DS_Store"}, ignoreLines)

	noIgnoreLines := readIgnoreFile(".gitignore")
	assert.Equal(t, []string{}, noIgnoreLines)
}
