package rule

// Options are options that can be configured and applied on a per-rule basis
type Options struct {
	WordBoundary bool  `yaml:"word_boundary"`
	IncludeNote  *bool `yaml:"include_note"`
}
