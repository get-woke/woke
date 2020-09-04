package parser

import (
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

// parseViolations returns all the violations (FileResults) for every valid file in each path
func (p *Parser) processViolations(paths []string) (fr []result.FileResults) {
	var wg sync.WaitGroup

	rchan := make(chan *result.FileResults)

	for i := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			ch := p.processViolationInPath(path)
			for r := range ch {
				rchan <- r
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

func (p *Parser) processViolationInPath(path string) chan *result.FileResults {
	files := p.walkDir(path)

	rchan := make(chan *result.FileResults)
	var wg sync.WaitGroup

	for f := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()

			v, err := generateFileViolationsFromFilename(f, p.Rules)
			if err != nil {
				log.Error().Err(err).Str("file", f).Msg("generateFileViolationsFromFilename error")
				return
			}

			if len(v.Results) == 0 {
				return
			}

			rchan <- v
		}(f)
	}
	go func() {
		wg.Wait()
		close(rchan)
	}()
	return rchan
}

func (p *Parser) walkDir(dirname string) <-chan string {
	paths := make(chan string)

	go func() {
		err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
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
			case <-time.After(time.Second * 30):
				return fmt.Errorf("walk timeout: %s", dirname)
			}
		})
		if err != nil {
			log.Error().Err(err).Str("dir", dirname).Msg("filepath.Walk error")
		}
		close(paths)
	}()

	return paths
}
