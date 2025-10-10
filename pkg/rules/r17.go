// Rule 17: Check for parameter overrides in helm directive
// This rule checks if applications have parameter overrides in helm configuration
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule17 struct {
}

func (r *Rule17) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R17",
		Name:            "Check for parameter overrides in helm directive",
		PassDescription: "Applications do not use parameter overrides in helm configuration.",
		FailDescription: "Applications use parameter overrides in helm configuration. Consider using valueFiles instead of hardcoded parameters for better maintainability.",
		Severity:        SeverityHigh,
	}
}

func (r *Rule17) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	// Iterate over all ArgoCD Application files
	for _, app := range file_analysis.Applications {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(app.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if helm directive exists
		_, hasHelm := getYAMLValue(&yamlDoc, "spec.source.helm")
		if !hasHelm {
			// No helm directive, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Check if helm directive contains parameters field
		_, hasParameters := getYAMLValue(&yamlDoc, "spec.source.helm.parameters")
		if hasParameters {
			// Has parameters field - rule fails
			failedFiles = append(failedFiles, app.Name)
		} else {
			// Has helm directive but no parameters - rule passes
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
