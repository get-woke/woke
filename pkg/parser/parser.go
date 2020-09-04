package parser

import (
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/get-woke/woke/pkg/util"
	"github.com/get-woke/woke/pkg/walker"
	"github.com/rs/zerolog/log"
)

var DefaultPath = []string{"."}

// TODO: can this be dynamically determined?
const numWorkers = 20

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
			log.Error().Err(r.err).Msg("violations error")
		}

		sort.Sort(r.res)
		fr = append(fr, r.res)
	}
	return
}

func (p *Parser) processFiles(files <-chan string, rchan chan _result, done chan bool, wg *sync.WaitGroup) {
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
}

func (p *Parser) processViolationInPath(path string, done chan bool, rchan chan _result) {
	var wg sync.WaitGroup

	files := p.walkDir(path, done)

	// run parallel, but bounded
	numWorkerStr := util.GetEnvDefault("WORKER_POOL_COUNT", "0")
	if numWorker, err := strconv.Atoi(numWorkerStr); err == nil && numWorker > 0 {
		log.Debug().Str("path", path).Str("type", "bounded").Int("workers", numWorker).Msg("process files")

		wg.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			go func() {
				p.processFiles(files, rchan, done, &wg)
				wg.Done()
			}()
		}

	} else {
		// run parallel unbounded. Potential high memory consumption
		log.Debug().Str("path", path).Str("type", "parallel").Msg("process files")

		p.processFiles(files, rchan, done, &wg)
	}
	wg.Wait()
}

func (p *Parser) walkDir(dirname string, done chan bool) <-chan string {
	paths := make(chan string)

	go func() {
		defer close(paths)
		_ = walker.Walk(dirname, func(path string, _ os.FileMode) error {
			if p.Ignorer != nil && p.Ignorer.Match(path) {
				return nil
			}

			if util.IsTextFileFromFilename(path) != nil {
				return nil
			}

			paths <- path
			return nil
		})
	}()

	return paths
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
