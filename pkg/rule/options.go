package rule

// Options are options that can be configured and applied on a per-rule basis
type Options struct {
	WordBoundary      bool  `yaml:"word_boundary"`
	WordBoundaryStart bool  `yaml:"word_boundary_start"`
	WordBoundaryEnd   bool  `yaml:"word_boundary_end"`
	IncludeNote       *bool `yaml:"include_note"`
}
