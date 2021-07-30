package ignore

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

type DirTestSuite struct {
	suite.Suite
	GFS billy.Filesystem // git repository root
}

func (suite *DirTestSuite) SetupTest() {
	// setup generic git repository root
	fs := memfs.New()
	f, err := fs.Create(".gitignore")
	suite.NoError(err)
	err = fs.MkdirAll(".git", os.ModePerm)
	suite.NoError(err)
	_, err = f.Write([]byte("vendor/g*/\n"))
	suite.NoError(err)
	_, err = f.Write([]byte("ignore.crlf\r\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	err = fs.MkdirAll("vendor", os.ModePerm)
	suite.NoError(err)
	f, err = fs.Create("vendor/.gitignore")
	suite.NoError(err)
	_, err = f.Write([]byte("!github.com/\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	err = fs.MkdirAll("another", os.ModePerm)
	suite.NoError(err)
	err = fs.MkdirAll("ignore.crlf", os.ModePerm)
	suite.NoError(err)
	err = fs.MkdirAll("vendor/github.com", os.ModePerm)
	suite.NoError(err)
	err = fs.MkdirAll("vendor/gopkg.in", os.ModePerm)
	suite.NoError(err)

	err = fs.MkdirAll("noignore", os.ModePerm)
	suite.NoError(err)
	f, err = fs.Create("noignore/.gitignore")
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	suite.GFS = fs
}

func (suite *DirTestSuite) TestReadIgnoreFile() {
	ignoreLines, _ := readIgnoreFile(suite.GFS, []string{}, ".gitignore")
	patterns := []gitignore.Pattern{
		gitignore.ParsePattern("vendor/g*/", []string{}),
		gitignore.ParsePattern("ignore.crlf", []string{}),
	}
	suite.Equal(patterns, ignoreLines)

	noIgnoreLines, _ := readIgnoreFile(suite.GFS, []string{"noignore"}, ".gitignore")
	suite.Nil(noIgnoreLines)
}

func (suite *DirTestSuite) TestReadPatterns() {
	ps, err := readPatterns(suite.GFS, nil)
	suite.Nil(err)
	suite.Len(ps, 3)

	m := gitignore.NewMatcher(ps)
	suite.True(m.Match([]string{"ignore.crlf"}, true))
	suite.True(m.Match([]string{"vendor", "gopkg.in"}, true))
	suite.False(m.Match([]string{"vendor", "github.com"}, true))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDirTestSuite(t *testing.T) {
	suite.Run(t, new(DirTestSuite))
}
