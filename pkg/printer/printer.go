package printer

import "github.com/caitlinelfring/woke/pkg/result"

// Printer is an interface for printing FileResults
type Printer interface {
	Print(*result.FileResults) error
}
