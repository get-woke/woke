package result

import "strings"

// FileResults contains all the Results for the file
type FileResults struct {
	Filename string
	Results  []Result
}

func (fr *FileResults) String() string {
	lines := []string{fr.Filename}
	for _, r := range fr.Results {
		lines = append(lines, r.String())
	}
	return strings.Join(lines, "\n")
}

// Len is part of sort.Interface
func (fr FileResults) Len() int {
	return len(fr.Results)
}

// Swap is part of sort.Interface
func (fr FileResults) Swap(i, j int) {
	fr.Results[i], fr.Results[j] = fr.Results[j], fr.Results[i]
}

// Less is part of sort.Interface
func (fr FileResults) Less(i, j int) bool {
	if fr.Results[i].GetStartPosition().Line == fr.Results[j].GetStartPosition().Line {
		return fr.Results[i].GetStartPosition().Column < fr.Results[j].GetStartPosition().Column
	}

	return fr.Results[i].GetStartPosition().Line < fr.Results[j].GetStartPosition().Line
}
