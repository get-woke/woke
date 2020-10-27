package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTextFile(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		f, _ := os.Open("testdata/empty.txt")
		defer f.Close()
		err := IsTextFile(f)
		assert.EqualError(t, err, ErrFileEmpty.Error())
	})

	t.Run("html file", func(t *testing.T) {
		f, _ := os.Open("testdata/index.html")
		defer f.Close()
		err := IsTextFile(f)
		assert.NoError(t, err)
	})

	t.Run("binary file", func(t *testing.T) {
		f, _ := os.Open("testdata/binary.dat")
		defer f.Close()
		err := IsTextFile(f)
		assert.EqualError(t, err, ErrFileNotText.Error())
	})

	t.Run("text file", func(t *testing.T) {
		f, _ := os.Open("testdata/text.txt")
		defer f.Close()
		err := IsTextFile(f)
		assert.NoError(t, err)
	})

	t.Run("xml file", func(t *testing.T) {
		f, _ := os.Open("testdata/index.xml")
		defer f.Close()
		err := IsTextFile(f)
		assert.NoError(t, err)
	})

	t.Run("missing file", func(t *testing.T) {
		f, _ := os.Open("testdata/missing.txt")
		defer f.Close()
		err := IsTextFile(f)
		assert.Error(t, err)
	})

	t.Run("directory", func(t *testing.T) {
		f, _ := os.Open("testdata")
		defer f.Close()
		err := IsTextFile(f)
		assert.EqualError(t, err, ErrIsDir.Error())
	})
}

func TestIsTextFileFromFilename(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		err := IsTextFileFromFilename("testdata/empty.txt")
		assert.EqualError(t, err, ErrFileEmpty.Error())
	})

	t.Run("html file", func(t *testing.T) {
		err := IsTextFileFromFilename("testdata/index.html")
		assert.NoError(t, err)
	})

	t.Run("binary file", func(t *testing.T) {
		err := IsTextFileFromFilename("testdata/binary.dat")
		assert.EqualError(t, err, ErrFileNotText.Error())
	})

	t.Run("text file", func(t *testing.T) {
		err := IsTextFileFromFilename("testdata/text.txt")
		assert.NoError(t, err)
	})

	t.Run("xml file", func(t *testing.T) {
		err := IsTextFileFromFilename("testdata/index.xml")
		assert.NoError(t, err)
	})

	t.Run("missing file", func(t *testing.T) {
		err := IsTextFileFromFilename("testdata/missing.txt")
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("directory", func(t *testing.T) {
		err := IsTextFileFromFilename("testdata")
		assert.EqualError(t, err, ErrIsDir.Error())
	})
}
