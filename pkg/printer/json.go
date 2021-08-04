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
	writer  io.Writer
	newList bool
}

// NewJSON returns a new JSON printer
func NewJSON(w io.Writer) *JSON {
	return &JSON{writer: w, newList: true}
}

func (p *JSON) Start() error {
	fmt.Fprint(p.writer, `{"findings": [`)
	return nil
}

func (p *JSON) End() error {
	fmt.Fprint(p.writer, `]}`+"\n")
	return nil
}

// Print prints in FileResults as json
// NOTE: The JSON printer will bring each line result as a JSON string.
// It will not be presented as an array of FileResults. You will neeed to
// Split by new line to parse the full output
func (p *JSON) Print(fs *result.FileResults) error {
	var buf bytes.Buffer
	if p.newList == true {
		p.newList = false
	} else {
		fmt.Fprintf(p.writer, `,`) // Add comma between issues
	}
	err := json.NewEncoder(&buf).Encode(fs)
	fmt.Fprint(p.writer, buf.String()) // json Encoder already puts a new line in, so no need for Println here
	return err
}
