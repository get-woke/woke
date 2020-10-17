package printer

import (
	"encoding/json"
	"fmt"

	"github.com/get-woke/woke/pkg/result"
)

// JSON is a JSON printer meant for a machine to read
type JSON struct{}

// NewJSON returns a new JSON printer
func NewJSON() *JSON {
	return &JSON{}
}

// Print prints in FileResults as json
// NOTE: The JSON printer will bring each line result as a JSON string.
// It will not be presented as an array of FileResults. You will neeed to
// Split by new line to parse the full output
func (p *JSON) Print(fs *result.FileResults) error {
	b, err := json.Marshal(fs)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
