package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/get-woke/woke/pkg/util"
	"github.com/rs/zerolog/log"
)

var DefaultPath = []string{"."}

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
		paths = DefaultPath
	}

	return p.processViolations(paths), nil
}

type _result struct {
	res result.FileResults
	err error
}

// parseViolations returns all the violations (FileResults) for every valid file in each path
func (p *Parser) processViolations(paths []string) (fr []result.FileResults) {
	var wg sync.WaitGroup

	rchan := make(chan _result)
	done := make(chan bool)
	defer close(done)

	for i := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			p.processViolationInPath(path, done, rchan)
		}(paths[i])
	}

	go func() {
		wg.Wait()
		close(rchan)
	}()

	for r := range rchan {
		if r.err != nil {
			log.Error().Err(r.err).Msg("filepath.Walk error")
		}

		sort.Sort(r.res)
		fr = append(fr, r.res)
	}
	return
}

func (p *Parser) processViolationInPath(path string, done chan bool, rchan chan _result) {
	var wg sync.WaitGroup

	files, errc := p.walkDir(path, done)
	for f := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()

			v, err := generateFileViolationsFromFilename(f, p.Rules)
			if len(v.Results) == 0 {
				return
			}
			select {
			case rchan <- _result{*v, err}:
			case <-done:
				return
			}
		}(f)
	}

	wg.Wait()

	if err := <-errc; err != nil {
		rchan <- _result{result.FileResults{}, err}
	}
}

func (p *Parser) walkDir(dirname string, done chan bool) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
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
				return nil
			case <-done:
				return errors.New("walk canceled")
			case <-time.After(time.Second * 30):
				return fmt.Errorf("walk timeout: %s", dirname)
			}
		})
	}()

	return paths, errc
}

// pathsIncludeStdin returns true if any element of the slice is stdin
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
