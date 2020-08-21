package parser

import (
	"bufio"
	"go/token"
	"os"

	"github.com/caitlinelfring/woke/pkg/rule"
)

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

	var results []*rule.Result

	// Splits on newlines by default.
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
