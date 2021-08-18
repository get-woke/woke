package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/get-woke/woke/pkg/result"
)

// JSON is a JSON printer meant for a machine to read
type JSON struct {
	writer io.Writer
}

// NewJSON returns a new JSON printer
func NewJSON(w io.Writer) *JSON {
	return &JSON{writer: w}
}

func (p *JSON) PrintSuccessExitMessage() bool {
	return true
}

func (p *JSON) Start() {
}

func (p *JSON) End() {
}

// Print prints in FileResults as json
// NOTE: The JSON printer will bring each line result as a JSON string.
// It will not be presented as an array of FileResults. You will neeed to
// Split by new line to parse the full output
func (p *JSON) Print(fs *result.FileResults) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(fs)
	fmt.Fprint(p.writer, buf.String()) // json Encoder already puts a new line in, so no need for Println here
	return err
}
