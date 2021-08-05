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
	writer     io.Writer
	newList    bool
	isTrueJSON bool
}

// NewJSON returns a new JSON printer
func NewJSON(w io.Writer, isTrueJSON bool) *JSON {
	return &JSON{writer: w, newList: true, isTrueJSON: isTrueJSON}
}

func (p *JSON) ShouldSkipExitMessage() bool {
	return p.isTrueJSON
}

func (p *JSON) Start() {
	if p.isTrueJSON {
		fmt.Fprint(p.writer, `{"findings": [`)
	}
}

func (p *JSON) End() {
	if p.isTrueJSON {
		fmt.Fprint(p.writer, `]}`+"\n")
	}
}

// Print prints in FileResults as json
// NOTE: The JSON printer will bring each line result as a JSON string.
// It will not be presented as an array of FileResults. You will neeed to
// Split by new line to parse the full output
func (p *JSON) Print(fs *result.FileResults) error {
	var buf bytes.Buffer
	if p.newList {
		p.newList = false
	} else if p.isTrueJSON {
		_, err := fmt.Fprint(p.writer, `,`) // Add comma between issues
		if err != nil {
			return err
		}
	}
	err := json.NewEncoder(&buf).Encode(fs)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(p.writer, buf.String()) // json Encoder already puts a new line in, so no need for Println here
	return err
}
