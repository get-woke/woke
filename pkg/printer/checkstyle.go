package printer

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/get-woke/woke/pkg/result"
)

// Checkstyle is a Checkstyle printer meant for use by a Checkstyle annotation
type Checkstyle struct {
	writer  io.Writer
	encoder *xml.Encoder
}

// NewCheckstyle returns a new Checkstyle printer
func NewCheckstyle(w io.Writer) *Checkstyle {
	return &Checkstyle{
		writer:  w,
		encoder: xml.NewEncoder(w),
	}
}

type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Errors  []Error  `xml:"error"`
}
type Error struct {
	XMLName  xml.Name `xml:"error"`
	Column   int      `xml:"column,attr"`
	Line     int      `xml:"line,attr"`
	Message  string   `xml:"message,attr"`
	Severity string   `xml:"severity,attr"`
	Source   string   `xml:"source,attr"`
}

func (p *Checkstyle) PrintSuccessExitMessage() bool {
	return true
}

// Print prints in the format for Checkstyle.
// https://github.com/checkstyle/checkstyle
func (p *Checkstyle) Print(fs *result.FileResults) error {
	var f File
	f.Name = fs.Filename
	for _, r := range fs.Results {
		f.Errors = append(f.Errors, Error{
			Column:   r.GetStartPosition().Column,
			Line:     r.GetStartPosition().Line,
			Message:  r.Reason(),
			Severity: r.GetSeverity().String(),
			Source:   "woke",
		})
	}
	return p.encoder.Encode(f)
}

func (p *Checkstyle) Start() {
	fmt.Fprint(p.writer, xml.Header)
	p.encoder.Indent("", "  ")
	p.encoder.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "checkstyle"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "version"}, Value: "5.0"},
		},
	})
	p.encoder.Flush()
}

func (p *Checkstyle) End() {
	p.encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "checkstyle"}})
	p.encoder.Flush()
}
