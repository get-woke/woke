package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTextFile(t *testing.T) {
	f1, _ := os.Open("testdata/empty.txt")
	defer f1.Close()
	err := IsTextFile(f1)
	assert.EqualError(t, err, ErrFileEmpty.Error())

	f2, _ := os.Open("testdata/binary.dat")
	defer f2.Close()
	err2 := IsTextFile(f2)
	assert.EqualError(t, err2, ErrFileNotText.Error())

	f3, _ := os.Open("testdata/text.txt")
	defer f3.Close()
	err3 := IsTextFile(f3)
	assert.NoError(t, err3)

	f4, _ := os.Open("testdata/missing.txt")
	defer f4.Close()
	err4 := IsTextFile(f4)
	assert.Error(t, err4)

	f5, _ := os.Open("testdata")
	defer f5.Close()
	err5 := IsTextFile(f5)
	assert.EqualError(t, err5, ErrIsDir.Error())
}

func TestIsTextFileFromFilename(t *testing.T) {
	err := IsTextFileFromFilename("testdata/empty.txt")
	assert.EqualError(t, err, ErrFileEmpty.Error())

	err2 := IsTextFileFromFilename("testdata/binary.dat")
	assert.EqualError(t, err2, ErrFileNotText.Error())

	err3 := IsTextFileFromFilename("testdata/text.txt")
	assert.NoError(t, err3)

	err4 := IsTextFileFromFilename("testdata/missing.txt")
	assert.True(t, os.IsNotExist(err4))

	err5 := IsTextFileFromFilename("testdata")
	assert.EqualError(t, err5, ErrIsDir.Error())
}
