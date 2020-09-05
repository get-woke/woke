package rule

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRules(t *testing.T) {
	for _, r := range DefaultRules {
		for _, term := range r.Terms {
			t.Run(r.Name+"/"+term, func(t *testing.T) {
				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("%s with other words after", term)), 1)
				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("other words before %s", term)), 1)
				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("other words %s before", term)), 1)
				assert.Len(t, r.FindAllStringIndex(term), 1)

				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("%s with other words after %s", term, term)), 2)
				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("%s other %s words before", term, term)), 2)
				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("other %s words %s before", term, term)), 2)
				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("other %s words before %s", term, term)), 2)

				assert.Len(t, r.FindAllStringIndex(fmt.Sprintf("other %s%s.", term, term)), 0)
			})
		}
	}
}
