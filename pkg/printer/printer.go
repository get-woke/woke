package printer

import "github.com/caitlinelfring/woke/pkg/result"

type Printer interface {
	Print(*result.FileResults) error
}
