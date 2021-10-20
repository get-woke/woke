package ignore

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

type IgnoreTestSuite struct {
	suite.Suite
	GFS billy.Filesystem // git repository root
}

func (suite *IgnoreTestSuite) TempFileSystem() (fs billy.Filesystem, clean func()) {
	fs = osfs.New(os.TempDir())
	path, err := util.TempDir(fs, "", "")
	if err != nil {
		panic(err)
	}

	fs, err = fs.Chroot(path)
	if err != nil {
		panic(err)
	}

	return fs, func() {
		_ = util.RemoveAll(fs, path)
	}
}

func (suite *IgnoreTestSuite) SetupTest() {
	// setup generic git repository root
	fs, clean := suite.TempFileSystem()
	defer clean()
	f, err := fs.Create(".gitignore")
	suite.NoError(err)
	_, err = f.Write([]byte("*.DS_Store\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	err = fs.MkdirAll(".git", os.ModePerm)
	suite.NoError(err)

	f, err = fs.Create(".ignore")
	suite.NoError(err)
	_, err = f.Write([]byte("*.IGNORE\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	f, err = fs.Create(".notignored")
	suite.NoError(err)
	_, err = f.Write([]byte("*.NOTIGNORED\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	f, err = fs.Create(".wokeignore")
	suite.NoError(err)
	_, err = f.Write([]byte("*.WOKEIGNORE\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	suite.GFS = fs
	// suite.tempDir = path
}

func (suite *IgnoreTestSuite) TearDownSuite() {
	// _ = util.RemoveAll(suite.GFS, suite.tempDir)
}

func (suite *IgnoreTestSuite) TestGetDomainFromWorkingDir() {
	suite.Equal([]string{}, getDomainFromWorkingDir("a/b/c/d", "b/c/d"))
	// suite.Equal([]string{}, getDomainFromWorkingDir("a/b/c/d", "a/b/c/d"))
	suite.Equal([]string{"d"}, getDomainFromWorkingDir("a/b/c/d", "c"))
	suite.Equal([]string{"d"}, getDomainFromWorkingDir("a/b/c/d", "b/c"))
	suite.Equal([]string{"b", "c", "d"}, getDomainFromWorkingDir("a/b/c/d", "a"))
	suite.Equal([]string{"c", "d"}, getDomainFromWorkingDir("a/b/c/d", "b/"))
	suite.Equal([]string{"b", "c", "d"}, getDomainFromWorkingDir("b/b/c/d", "b/"))
}

func (suite *IgnoreTestSuite) TestGetRootGitDir() {
	fs, clean := suite.TempFileSystem()
	log.Println(fs.Root())
	defer clean()
	dir, err := util.TempDir(fs, fs.Root(), "")
	log.Println(filepath.Abs(dir))
	suite.NoError(err)

	subdir := filepath.Join(dir, "a")
	log.Println(filepath.Abs(subdir))
	err = fs.MkdirAll(subdir, 0755)
	suite.NoError(err)

	gitdir := filepath.Join(dir, ".git")
	log.Println(filepath.Abs(gitdir))
	err = fs.MkdirAll(gitdir, 0755)
	suite.NoError(err)

	// root := fs.Root()

	// expected := osfs.New(fs.Join(root, dir))
	rootFs, err := GetRootGitDir(subdir)
	suite.NoError(err)
	suite.Equal(dir, rootFs.Root())
}
func (suite *IgnoreTestSuite) TestIgnore_Match() {
	i, err := NewIgnore(suite.GFS, []string{"my/files/*"})
	suite.NoError(err)
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
	i, err := NewIgnore(suite.GFS, []string{"*.FROMARGUMENT"})
	suite.NoError(err)
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
