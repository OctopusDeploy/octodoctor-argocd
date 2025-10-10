// Rule 15: Check for hardcoded values in helm directive
// This rule checks if applications have hardcoded values in helm configuration
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule15 struct {
}

func (r *Rule15) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R15",
		Name:            "Check for hardcoded values in helm directive",
		PassDescription: "Applications do not use hardcoded values in helm configuration.",
		FailDescription: "Applications use hardcoded values in helm configuration. Consider using valueFiles instead of hardcoded values for better maintainability.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule15) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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
		helmNode, hasHelm := getYAMLValue(&yamlDoc, "spec.source.helm")
		if !hasHelm {
			// No helm directive, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Check if helm directive contains "values" or "valuesObject" fields
		hasValues := false
		hasValuesObject := false

		if helmNode.Kind == yaml.MappingNode {
			// Iterate through the helm mapping to check for "values" or "valuesObject"
			for i := 0; i < len(helmNode.Content); i += 2 {
				if helmNode.Content[i].Value == "values" {
					hasValues = true
				}
				if helmNode.Content[i].Value == "valuesObject" {
					hasValuesObject = true
				}
			}
		}

		if hasValues || hasValuesObject {
			// Has values or valuesObject field - rule fails
			failedFiles = append(failedFiles, app.Name)
		} else {
			// Has helm directive but no values/valuesObject - rule passes
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
