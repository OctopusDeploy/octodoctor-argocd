package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// renderJSON writes the JSON report of result to a temporary directory and returns its content
func renderJSON(t *testing.T, result *analysis.ArgoCDAnalysis) []byte {
	t.Helper()
	dir := t.TempDir()
	GenerateStructuredReport(result, dir, "json", "test")
	data, err := os.ReadFile(filepath.Join(dir, "report.json"))
	require.NoError(t, err)
	return data
}

func TestStructuredReportIgnored(t *testing.T) {
	t.Run("Ignored rules", func(t *testing.T) {
		ignored := analysis.IgnoredExecution{
			RuleExecution: analysis.RuleExecution{ID: "R19", FilesTriggered: []string{"applications/mesh-a.yaml"}},
			Reason:        "pinned",
		}
		var report JSONReport
		require.NoError(t, json.Unmarshal(renderJSON(t, &analysis.ArgoCDAnalysis{
			RulesFailed:       []analysis.RuleExecution{{ID: "R14"}},
			RulesIgnored:      []analysis.IgnoredExecution{ignored},
			RulesFullyIgnored: 1,
		}), &report))

		assert.Equal(t, Rules{Total: 2, Failed: 1, Ignored: 1}, report.Summary.Rules)
		assert.Equal(t, []Rule{{Id: "R19", AffectedFiles: []string{"applications/mesh-a.yaml"}, Reason: "pinned"}}, report.Results.Ignored)
	})

	t.Run("No ignored rules", func(t *testing.T) {
		data := string(renderJSON(t, &analysis.ArgoCDAnalysis{}))
		assert.Contains(t, data, `"ignored": []`)
		assert.NotContains(t, data, `"reason"`)
	})
}
