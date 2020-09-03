package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/get-woke/woke/pkg/util"
)

// Parser parses files and finds lines that break rules
type Parser struct {
	Rules   []*rule.Rule
	Ignorer *ignore.Ignore
}

func NewParser(rules []*rule.Rule, ignorer *ignore.Ignore) *Parser {
	return &Parser{
		Rules:   rules,
		Ignorer: ignorer,
	}
}

// ParsePaths parses all files provided and returns the results
func (p *Parser) ParsePaths(paths ...string) (results []result.FileResults, err error) {
	// data provided through stdin
	if pathsIncludeStdin(paths) {
		r, err := generateFileViolations(os.Stdin, p.Rules)
		return []result.FileResults{*r}, err
	}

	if len(paths) == 0 {
		paths = []string{"."}
	}

	return p.processViolations(paths)
}

func pathsIncludeStdin(paths []string) bool {
	if len(paths) == 0 {
		return false
	}
	for _, p := range paths {
		if p == os.Stdin.Name() {
			return true
		}
	}
	return false
}

func (p *Parser) processViolations(paths []string) (fr []result.FileResults, err error) {
	var wg sync.WaitGroup

	done := make(chan struct{})
	defer close(done)
	rchan := make(chan *result.FileResults)

	for i := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			innerrchan := p.processViolationInPath(path, done)
			for r := range innerrchan {
				select {
				case rchan <- r:
				case <-done:
				}
			}
		}(paths[i])
	}

	go func() {
		wg.Wait()
		close(rchan)
	}()

	for r := range rchan {
		sort.Sort(r)
		fr = append(fr, *r)
	}
	return
}

func (p *Parser) processViolationInPath(path string, done <-chan struct{}) chan *result.FileResults {
	// TODO: Handler errors
	files, _ := p.walkDir(path, done)
	rchan := make(chan *result.FileResults)

	var wg sync.WaitGroup
	for f := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()

			v, err := generateFileViolationsFromFilename(f, p.Rules)
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(v.Results) == 0 {
				return
			}

			select {
			case rchan <- v:
			case <-done:
			}
		}(f)
	}
	go func() {
		wg.Wait()
		close(rchan)
	}()
	return rchan
}

func (p *Parser) walkDir(dirname string, done <-chan struct{}) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)

		errc <- filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if p.Ignorer != nil && p.Ignorer.Match(path) {
				return nil
			}
			if util.IsTextFileFromFilename(path) != nil {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}
