// Rule 22: Check argocd-secret webhook configuration
// This rule checks if webhook secrets are configured for Git providers
package rules

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule22 struct {
}

func (r *Rule22) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R22",
		Name:            "Check argocd-secret webhook configuration",
		PassDescription: "Webhooks are configured in argocd-secret for improved performance.",
		FailDescription: "No webhooks are configured in argocd-secret. Webhooks improve Argo CD performance by triggering immediate syncs instead of polling.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule22) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	foundArgoCdSecret := false

	// Iterate over all ArgoCD Secret files
	for _, secret := range file_analysis.Secrets {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(secret.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if this is the argocd-secret Secret
		nameNode, hasName := getYAMLValue(&yamlDoc, "metadata.name")
		if !hasName || nameNode.Kind != yaml.ScalarNode || nameNode.Value != "argocd-secret" {
			// Not the argocd-secret Secret, skip this file
			inactiveFiles = append(inactiveFiles, secret.Name)
			continue
		}

		foundArgoCdSecret = true

		// Check if any webhook properties exist in data
		dataNode, hasData := getYAMLValue(&yamlDoc, "data")
		hasWebhook := false

		if hasData && dataNode.Kind == yaml.MappingNode {
			// Look for any key starting with "webhook" in the data mapping
			for i := 0; i < len(dataNode.Content); i += 2 {
				keyNode := dataNode.Content[i]
				if strings.HasPrefix(keyNode.Value, "webhook") {
					hasWebhook = true
					break
				}
			}
		}

		if !hasWebhook {
			// No webhook properties found, rule fails
			failedFiles = append(failedFiles, secret.Name)
		} else {
			// At least one webhook property exists, rule passes
			passedFiles = append(passedFiles, secret.Name)
		}
	}

	// If no argocd-secret Secret was found, mark the rule as inactive
	if !foundArgoCdSecret {
		result := analysis.RuleExecution{
			ID:             r.Metadata().ID,
			Name:           r.Metadata().Name,
			Description:    r.Metadata().PassDescription,
			Severity:       r.Metadata().Severity,
			FilesTriggered: nil,
		}
		return RuleResult{
			Active:    false,
			Passed:    true,
			Execution: result,
		}, nil
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
