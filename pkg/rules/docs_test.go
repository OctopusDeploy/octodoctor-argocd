package rules

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEveryRuleHasDocumentation guards the DocSlug of every registered rule against the
// docs tree. The published URLs cannot be verified from here, so the next best check is
// that each slug names a page that actually exists under docs/rules.
func TestEveryRuleHasDocumentation(t *testing.T) {
	seen := make(map[string]string)

	for _, rule := range NewRegistry().Rules {
		metadata := rule.Metadata()

		require.NotEmpty(t, metadata.DocSlug, "rule %s has no DocSlug", metadata.ID)
		assert.FileExists(t, filepath.Join("..", "..", "docs", "rules", metadata.DocSlug+".md"),
			"rule %s points at a documentation page that does not exist", metadata.ID)

		other, duplicate := seen[metadata.DocSlug]
		assert.False(t, duplicate, "rules %s and %s share the DocSlug %q", other, metadata.ID, metadata.DocSlug)
		seen[metadata.DocSlug] = metadata.ID

		assert.Equal(t,
			"https://octodoctor-argocd.readthedocs.io/en/latest/rules/"+metadata.DocSlug+"/",
			metadata.DocumentationURL(),
			"rule %s builds an unexpected documentation URL", metadata.ID)
	}
}
