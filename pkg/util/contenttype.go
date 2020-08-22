package util

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func detectContentType(file *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	// Reset the file so a scanner can scan
	_, _ = file.Seek(0, 0)

	if err != nil {
		return "", err
	}
	return http.DetectContentType(buffer[:n]), nil
}

func isTextFile(file *os.File) bool {
	contentType, err := detectContentType(file)
	if err != nil {
		return false
	}

	return strings.HasPrefix(contentType, "text/plain")
}

func IsTextFile(file *os.File) error {
	if !isTextFile(file) {
		return fmt.Errorf("%s is not a text file", file.Name())
	}
	return nil
}
