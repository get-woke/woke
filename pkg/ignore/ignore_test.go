package ignore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

type IgnoreTestSuite struct {
	suite.Suite
	GFS billy.Filesystem // git repository root
}

func TempFileSystem() (fs billy.Filesystem, clean func()) {
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
	fs, clean := TempFileSystem()
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
	_, err = f.Write([]byte("*.WOKEIGNORE\ntestFolder\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	err = fs.MkdirAll("nestedIgnoreFolder", os.ModePerm)
	suite.NoError(err)

	f, err = fs.Create(fs.Join("nestedIgnoreFolder", ".wokeignore"))
	suite.NoError(err)
	_, err = f.Write([]byte("*.NESTEDIGNORE\n"))
	suite.NoError(err)
	err = f.Close()
	suite.NoError(err)

	suite.GFS = fs
}

func BenchmarkIgnore(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
	fs, clean := TempFileSystem()
	defer clean()
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			err := fs.MkdirAll(fs.Join(fmt.Sprintf("%d", i), fmt.Sprintf("%d", j)), os.ModePerm)
			assert.NoError(b, err)
			f, err := fs.Create(".wokeignore")
			assert.NoError(b, err)
			_, err = f.Write([]byte("testFolder"))
			assert.NoError(b, err)
			err = f.Close()
			assert.NoError(b, err)
			for k := 0; k < 100; k++ {
				f, err := fs.Create(fmt.Sprintf("%d.txt", k))
				assert.NoError(b, err)
				err = f.Close()
				assert.NoError(b, err)
			}
		}
	}

	b.ResetTimer()
	b.Run("nested-ignores-enabled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ignorer, err := NewIgnore(fs, []string{}, false)
			assert.NoError(b, err)
			ignorer.Match(filepath.Join("not", "foo"), false)
		}
	})
	b.Run("nested-ignores-disabled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ignorer, err := NewIgnore(fs, []string{}, true)
			assert.NoError(b, err)
			ignorer.Match(filepath.Join("not", "foo"), false)
		}
	})
}

func (suite *IgnoreTestSuite) TestGetDomainFromWorkingDir() {
	suite.Equal([]string{}, getDomainFromWorkingDir(filepath.FromSlash("a/b/c/d"), filepath.FromSlash("b/c/d")))
	suite.Equal([]string{}, getDomainFromWorkingDir(filepath.FromSlash("a/b/c/d"), filepath.FromSlash("a/b/c/d")))
	suite.Equal([]string{"d"}, getDomainFromWorkingDir(filepath.FromSlash("a/b/c/d"), "c"))
	suite.Equal([]string{"d"}, getDomainFromWorkingDir(filepath.FromSlash("a/b/c/d"), filepath.FromSlash("b/c")))
	suite.Equal([]string{"b", "c", "d"}, getDomainFromWorkingDir(filepath.FromSlash("a/b/c/d"), "a"))
	suite.Equal([]string{"c", "d"}, getDomainFromWorkingDir(filepath.FromSlash("a/b/c/d"), filepath.FromSlash("b/")))
	suite.Equal([]string{"b", "c", "d"}, getDomainFromWorkingDir(filepath.FromSlash("b/b/c/d"), filepath.FromSlash("b/")))
}

func (suite *IgnoreTestSuite) TestGetRootGitDir() {
	cwd, err := os.Getwd()
	suite.NoError(err)

	rootFs, err := GetRootGitDir(cwd)
	suite.NoError(err)
	suite.Equal(filepath.Dir(filepath.Dir(cwd)), rootFs.Root())
}

func (suite *IgnoreTestSuite) TestGetRootGitDirNotExist() {
	fs, clean := TempFileSystem()
	defer clean()
	rootFs, err := GetRootGitDir(fs.Root())
	suite.NoError(err)
	suite.Equal(fs.Root(), rootFs.Root())
}

func (suite *IgnoreTestSuite) TestIgnoreLines_Match() {
	i, err := NewIgnore(suite.GFS, []string{"my/files/*"}, false)
	suite.NoError(err)
	suite.NotNil(i)

	suite.False(i.Match(filepath.Join("not", "foo"), false))
	suite.True(i.Match(filepath.Join("my", "files", "file1"), false))
	suite.False(i.Match(filepath.Join("my", "files"), false))
}

func (suite *IgnoreTestSuite) TestIgnoreLinesNoTraversal_Match() {
	i, err := NewIgnore(suite.GFS, []string{"my/files/*"}, true)
	suite.NoError(err)
	suite.NotNil(i)

	suite.False(i.Match(filepath.Join("not", "foo"), false))
	suite.True(i.Match(filepath.Join("my", "files", "file1"), false))
	suite.False(i.Match(filepath.Join("my", "files"), false))
}

// Test all default ignore files, except for .git/info/exclude, since
// that uses a .git directory that we cannot check in.
func (suite *IgnoreTestSuite) TestIgnoreDefaultIgnoreFiles_Match() {
	i, err := NewIgnore(suite.GFS, []string{"*.FROMARGUMENT"}, false)
	suite.NoError(err)
	suite.NotNil(i)

	// Test top-level ignore files all match
	suite.testCommonIgnoreDefaultIgnoreFilesMatch(i)

	// Test match from the nested ./nestedIgnoreFolder/.wokeignore when nested ignores is enabled
	suite.True(i.Match(filepath.Join("nestedIgnoreFolder", "testdata", "test.NESTEDIGNORE"), false))
}

func (suite *IgnoreTestSuite) TestIgnoreDefaultIgnoreFilesNoTraversal_Match() {
	i, err := NewIgnore(suite.GFS, []string{"*.FROMARGUMENT"}, true)
	suite.NoError(err)
	suite.NotNil(i)

	// Test top-level ignore files all match
	suite.testCommonIgnoreDefaultIgnoreFilesMatch(i)

	// Test no match from the nested ./nestedIgnoreFolder/.wokeignore when nested ignores is disabled
	suite.False(i.Match(filepath.Join("nestedIgnoreFolder", "testdata", "test.NESTEDIGNORE"), false))
}

func (suite *IgnoreTestSuite) testCommonIgnoreDefaultIgnoreFilesMatch(i *Ignore) {
	suite.False(i.Match(filepath.Join("testdata", "notfoo"), false))
	suite.True(i.Match(filepath.Join("testdata", "test.FROMARGUMENT"), false)) // From .gitignore
	suite.True(i.Match(filepath.Join("testdata", "test.DS_Store"), false))     // From .gitignore
	suite.True(i.Match(filepath.Join("testdata", "test.IGNORE"), false))       // From .ignore
	suite.True(i.Match(filepath.Join("testdata", "test.WOKEIGNORE"), false))   // From .wokeignore
	suite.True(i.Match(filepath.Join("testdata", "testFolder"), true))         // From .wokeignore
	suite.False(i.Match(filepath.Join("testdata", "notTestFolder"), true))     // From .wokeignore
	suite.False(i.Match(filepath.Join("testdata", "test.NOTIGNORED"), false))  // From .notincluded - making sure only default are included
}

func (suite *IgnoreTestSuite) TestReadIgnoreFile() {
	ignoreLines := readIgnoreFile(suite.GFS, ".gitignore")
	suite.Equal([]string{"*.DS_Store"}, ignoreLines)

	noIgnoreLines := readIgnoreFile(suite.GFS, "missing.gitignore")
	suite.Equal([]string{}, noIgnoreLines)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestIgnoreTestSuite(t *testing.T) {
	suite.Run(t, new(IgnoreTestSuite))
}
