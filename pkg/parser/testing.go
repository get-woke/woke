package parser

import (
	"io/ioutil"
	"os"
	"testing"
)

// newFileWithPrefix creates a new file for testing with the given prefix. The file,
// and the directory that the file was created in will be removed at the completion
// of the test
func newFileWithPrefix(t *testing.T, prefix, text string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	b := []byte(text)
	_, err = tmpFile.Write(b)

	return tmpFile, err
}

// newFile creates a new file for testing. The file, and the directory that the file
// was created in will be removed at the completion of the test
func newFile(t *testing.T, text string) (*os.File, error) {
	return newFileWithPrefix(t, "woke-", text)
}
