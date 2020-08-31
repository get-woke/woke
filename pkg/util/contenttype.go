package util

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

var ErrFileEmpty = errors.New("file is empty")
var ErrFileNotText = errors.New("file is not a text file")

func detectContentType(file *os.File) string {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	// Reset the file so a scanner can scan
	_, _ = file.Seek(0, 0)

	return http.DetectContentType(buffer[:n])
}

func isTextFile(file *os.File) bool {
	contentType := detectContentType(file)

	return strings.HasPrefix(contentType, "text/plain")
}

// IsTextFile returns an error if the file is not of content-type 'text/plain'
func IsTextFile(file *os.File) error {
	e, err := file.Stat()
	if err != nil {
		return err
	}
	if e.Size() == 0 {
		return ErrFileEmpty
	}

	if !isTextFile(file) {
		return ErrFileNotText
	}

	return nil
}
