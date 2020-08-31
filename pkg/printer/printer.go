package printer

import "github.com/get-woke/woke/pkg/result"

// Printer is an interface for printing FileResults
type Printer interface {
	Print(*result.FileResults) error
}
