package parser

import (
	"bufio"
	"fmt"
	"go/token"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/caitlinelfring/woke/pkg/rule"
)

// Parser parses files and finds lines that break rules
type Parser struct {
	Rules []*rule.Rule
}

func (p *Parser) ParseFiles(files []string) *rule.Results {
	results := rule.Results{}

	for _, f := range files {
		// skip rules config, which will always produce failures
		// if f == ruleConfig {
		// 	continue
		// }
		r, err := p.Parse(f)
		if err != nil {
			log.Printf("parser failed on %s: %s\n", f, err)
			continue
		}
		results.Push(r.Results...)
	}
	return &results
}

// Parse reads the file and returns results of places where rules are broken
func (p *Parser) Parse(filename string) (results rule.Results, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	if err = errIsNotTextFile(f); err != nil {
		return
	}

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
				rs := rule.Result{
					Rule:  r,
					Match: text[i[0]:i[1]],
					Position: &token.Position{
						Filename: filename,
						Line:     line,
						Column:   i[0],
					},
				}
				results.Add(&rs)
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

func errIsNotTextFile(file *os.File) error {
	if !isTextFile(file) {
		return fmt.Errorf("%s is not a text file", file.Name())
	}
	return nil
}
