package printer

import (
	"fmt"

	"github.com/caitlinelfring/woke/pkg/result"
	"github.com/caitlinelfring/woke/pkg/rule"
)

// GitHubActions is a GitHubActions printer meant for use by a GitHub Action annotation
type GitHubActions struct{}

// NewGitHubActions returns a new GitHubActions printer
func NewGitHubActions() *GitHubActions {
	return &GitHubActions{}
}

// Print prints in the format for GitHub actions
// https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-an-error-message
func (p *GitHubActions) Print(fs *result.FileResults) error {
	for _, r := range fs.Results {
		fmt.Printf("::%s file=%s,line=%d,col=%d::%s\n",
			translateSeverityForAction(r.Rule.Severity),
			r.StartPosition.Filename,
			r.StartPosition.Line,
			r.StartPosition.Column,
			r.Reason())
	}
	return nil
}

func translateSeverityForAction(s rule.Severity) string {
	if s == rule.SevError {
		return "error"
	}
	// treat everything else as a warning
	// https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions#setting-a-warning-message
	return "warning"
}
