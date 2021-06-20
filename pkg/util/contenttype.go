package util

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	// ErrFileEmpty is an error to signify the file is empty
	ErrFileEmpty = errors.New("file is empty")
	// ErrFileNotText is an error to signify the file type is not a text file that can be read
	ErrFileNotText = errors.New("file is not a text file")
	// ErrIsDir is an error to signify a directory
	ErrIsDir = errors.New("file is a directory")
)

func detectContentType(file io.Reader) string {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)

	return http.DetectContentType(buffer[:n])
}

func isTextFile(file *os.File) bool {
	contentType := detectContentType(file)

	return strings.HasPrefix(contentType, "text/")
}

// IsTextFileFromFilename returns an error if the filename is not of content-type 'text/*'
func IsTextFileFromFilename(filename string) error {
	// Don't check stdin to avoid closing it prematurely
	if filename == os.Stdin.Name() {
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return IsTextFile(f)
}

// IsTextFile returns an error if the file is not of content-type 'text/*'
func IsTextFile(file *os.File) error {
	e, err := file.Stat()
	if err != nil {
		return err
	}
	if e.IsDir() {
		return ErrIsDir
	}

	if e.Size() == 0 {
		return ErrFileEmpty
	}

	if !isTextFile(file) {
		return ErrFileNotText
	}

	return nil
}
