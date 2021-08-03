package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/get-woke/woke/pkg/result"
)

// SonarQube is a JSON printer meant for import into SonarQube
type SonarQube struct {
	writer  io.Writer
	newList bool
}

type TextRange struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndColumn   int `json:"endColumn"`
}

type Location struct {
	Message   string    `json:"message"`
	FilePath  string    `json:"filePath"`
	TextRange TextRange `json:"textRange"`
}

type Issue struct {
	EngineID        string   `json:"engineId"`
	RuleID          string   `json:"ruleId"`
	PrimaryLocation Location `json:"primaryLocation"`
	Type            string   `json:"type"`
	Severity        string   `json:"severity"`
}

// NewSonarQube returns a new SonarQube JSON printer
func NewSonarQube(w io.Writer) *SonarQube {
	return &SonarQube{writer: w, newList: true}
}

// Print prints in FileResults as json
// NOTE: The JSON printer will bring each line result as a JSON string.
// It will not be presented as an array of FileResults. You will neeed to
// Split by new line to parse the full output
func (p *SonarQube) Print(fs *result.FileResults) error {
	var issue Issue

	if p.newList == true {
		p.newList = false
	} else {
		fmt.Fprint(p.writer, `,`) // add comma between issues
	}

	for i, res := range fs.Results {

		if i != 0 {
			fmt.Fprint(p.writer, `,`) // add comma between issues in list
		}
		// Translate the severity to sonarqube terms
		var sev = `CRITICAL`
		if res.GetSeverity().String() == `warning` {
			sev = `MAJOR`
		} else if res.GetSeverity().String() == `info` {
			sev = `INFO`
		}

		issue = Issue{
			EngineID: `woke`,
			Type:     `CODE_SMELL`,
			Severity: sev,
			RuleID:   res.GetRuleName(),
			PrimaryLocation: Location{
				Message:  res.Reason(),
				FilePath: fs.Filename,
				TextRange: TextRange{
					StartLine:   res.GetStartPosition().Line,
					StartColumn: res.GetStartPosition().Column,
					EndColumn:   res.GetEndPosition().Column}}}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(issue)
		if err != nil {
			return err // bubble up errors instead of continuing
		}
		fmt.Fprint(p.writer, buf.String()) // json Encoder already puts a new line in, so no need for Println here
	}

	return nil
}

func (p *SonarQube) Start() error {
	fmt.Fprint(p.writer, `{"issues":[`)
	return nil
}

func (p *SonarQube) End() error {
	fmt.Fprint(p.writer, `]}`+"\n")
	return nil
}
