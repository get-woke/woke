package ignore

import (
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

type IgnoreTestSuite struct {
	suite.Suite
}

func (suite *IgnoreTestSuite) TestGetDomainFromWorkingDir() {
	suite.Equal([]string{}, getDomainFromWorkingDir("a/b/c/d", "b/c/d"))
	suite.Equal([]string{}, getDomainFromWorkingDir("a/b/c/d", "a/b/c/d"))
	suite.Equal([]string{"d"}, getDomainFromWorkingDir("a/b/c/d", "c"))
	suite.Equal([]string{"d"}, getDomainFromWorkingDir("a/b/c/d", "b/c"))
	suite.Equal([]string{"b", "c", "d"}, getDomainFromWorkingDir("a/b/c/d", "a"))
	suite.Equal([]string{"c", "d"}, getDomainFromWorkingDir("a/b/c/d", "b/"))
}

func (suite *IgnoreTestSuite) TestIgnore_Match() {
	i := NewIgnore([]string{"my/files/*"}, false)
	suite.NotNil(i)

	// Test if rules with backslashes match on windows
	suite.False(i.Match("not/foo", false))
	suite.True(i.Match("my/files/file1", false))
	suite.False(i.Match("my/files", false))

	suite.False(i.Match(filepath.Join("not", "foo"), false))
	suite.True(i.Match(filepath.Join("my", "files", "file1"), false))
	suite.False(i.Match(filepath.Join("my", "files"), false))
}

// Test all default ignore files, except for .git/info/exclude, since
// that uses a .git directory that we cannot check in.
func (suite *IgnoreTestSuite) TestIgnoreDefaultIgoreFiles_Match() {
	i := NewIgnore([]string{"*.FROMARGUMENT"}, false)
	suite.NotNil(i)

	suite.False(i.Match(filepath.Join("testdata", "notfoo"), false))
	suite.True(i.Match(filepath.Join("testdata", "test.FROMARGUMENT"), false)) // From .gitignore
	suite.True(i.Match(filepath.Join("testdata", "test.DS_Store"), false))     // From .gitignore
	suite.True(i.Match(filepath.Join("testdata", "test.IGNORE"), false))       // From .ignore
	suite.True(i.Match(filepath.Join("testdata", "test.WOKEIGNORE"), false))   // From .wokeignore
	suite.False(i.Match(filepath.Join("testdata", "test.NOTIGNORED"), false))  // From .notincluded - making sure only default are included
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestIgnoreTestSuite(t *testing.T) {
	suite.Run(t, new(IgnoreTestSuite))
}
