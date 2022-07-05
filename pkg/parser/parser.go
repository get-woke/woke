package parser

import (
	"os"
	"sort"
	"sync"

	"github.com/get-woke/woke/pkg/ignore"
	"github.com/get-woke/woke/pkg/printer"
	"github.com/get-woke/woke/pkg/result"
	"github.com/get-woke/woke/pkg/rule"
	"github.com/get-woke/woke/pkg/util"
	"github.com/get-woke/woke/pkg/walker"

	env "github.com/caitlinelfring/go-env-default"
	"github.com/rs/zerolog/log"
)

// DefaultPath is the default path if no paths are provided
var DefaultPath = []string{"."}

// TODO: can this be dynamically determined?
const numWorkers = 20

// Parser parses files and finds lines that break rules
type Parser struct {
	Rules   []*rule.Rule
	Ignorer *ignore.Ignore

	rchan chan result.FileResults
}

// NewParser returns a pointer to a Parser that is used to check for findings
// based on the rules provided, ignoring files based on the ignorer provided
func NewParser(rules []*rule.Rule, ignorer *ignore.Ignore) *Parser {
	return &Parser{
		Rules:   rules,
		Ignorer: ignorer,
		rchan:   make(chan result.FileResults),
	}
}

// ParsePaths parses all files provided and returns the number of files with findings
func (p *Parser) ParsePaths(print printer.Printer, paths ...string) int {
	print.Start()
	defer print.End()

	// data provided through stdin
	if util.InSlice(os.Stdin.Name(), paths) {
		r, _ := p.generateFileFindings(os.Stdin)
		if r.Len() > 0 {
			print.Print(r)
		}
		return r.Len()
	}

	var wg sync.WaitGroup

	done := make(chan bool)
	defer close(done)

	for _, path := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			p.processFindingInPath(path, done)
		}(path)
	}

	go func() {
		wg.Wait()
		close(p.rchan)
	}()

	findings := 0
	for r := range p.rchan {
		sort.Sort(r)
		print.Print(&r)
		findings++
	}
	return findings
}

func (p *Parser) processFiles(files <-chan string, done chan bool, wg *sync.WaitGroup) {
	for f := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()

			v, _ := p.generateFileFindingsFromFilename(f)
			if v == nil || len(v.Results) == 0 {
				return
			}
			p.rchan <- *v
		}(f)
	}
}

func (p *Parser) processFindingInPath(path string, done chan bool) {
	var wg sync.WaitGroup

	files := p.walkDir(path, done)

	// run parallel, but bounded
	numWorker := env.GetIntDefault("WORKER_POOL_COUNT", 0)
	if numWorker > 0 {
		log.Debug().Str("path", path).Str("type", "bounded").Int("workers", numWorker).Msg("process files")

		wg.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			go func() {
				p.processFiles(files, done, &wg)
				wg.Done()
			}()
		}
	} else {
		// run parallel unbounded. Potential high memory consumption
		log.Debug().Str("path", path).Str("type", "parallel").Msg("process files")

		p.processFiles(files, done, &wg)
	}

	wg.Wait()
}

func (p *Parser) walkDir(dirname string, done chan bool) <-chan string {
	paths := make(chan string)

	go func() {
		defer close(paths)
		_ = walker.Walk(dirname, func(path string, info os.FileMode) error {
			if p.Ignorer != nil && p.Ignorer.Match(path, info.IsDir()) {
				log.Debug().Str("file", path).Str("reason", "ignored file").Msg("skipping")
				return nil
			}

			paths <- path
			return nil
		})
	}()

	return paths
}
