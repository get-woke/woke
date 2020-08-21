package parser

import (
	"bufio"
	"go/token"
	"net/http"
	"os"
	"strings"

	"github.com/caitlinelfring/woke/pkg/rule"
)

// Parser parses files and finds lines that break rules
type Parser struct {
	Rules []*rule.Rule
}

// Parse reads the file and returns results of places where rules are broken
func (p *Parser) Parse(filename string) ([]*rule.Result, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	if !isTextFile(f) {
		// fmt.Printf("Skipping %s, is not a text file...\n", filename)
		return nil, nil
	}

	// Reset the file since we read it when checking the content-type
	_, _ = f.Seek(0, 0)

	var results []*rule.Result

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		text := scanner.Text()
		for _, r := range p.Rules {
			idx := r.Regexp.FindAllStringIndex(text, -1)
			if idx == nil {
				continue
			}

			for _, i := range idx {
				result := rule.Result{
					Rule:  r,
					Match: text[i[0]:i[1]],
					Position: &token.Position{
						Filename: filename,
						Line:     line,
						Column:   i[0],
					},
				}
				results = append(results, &result)
			}
		}

		line++
	}
	return results, scanner.Err()
}

func detectContentType(file *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
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
