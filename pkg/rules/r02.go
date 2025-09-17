// Rule 02: Check if sync policy is automated
// This rule checks if applications have automated sync policy configured
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule02 struct {
}

func (r *Rule02) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R02",
		Name:            "Check if sync policy is automated",
		PassDescription: "Applications have automated sync policy configured.",
		FailDescription: "Applications do not have automated sync policy configured. Manual sync is required for deployments.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule02) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string

	// Iterate over all ArgoCD Application files
	for _, app := range file_analysis.Applications {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(app.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if automated sync policy exists
		_, hasAutomated := getYAMLValue(&yamlDoc, "spec.syncPolicy.automated")
		if hasAutomated {
			// syncPolicy.automated exists, rule passes
			passedFiles = append(passedFiles, app.Name)
		} else {
			// syncPolicy.automated doesn't exist, rule fails
			failedFiles = append(failedFiles, app.Name)
		}
	}

	// Determine overall rule result
	// Note: This rule cannot be inactive - it either passes or fails for each application
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, nil)
}
