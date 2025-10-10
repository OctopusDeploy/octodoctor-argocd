// Rule 14: Check for prune setting in sync policy
// This rule checks if applications have prune=true in their sync policy
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule14 struct {
}

func (r *Rule14) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R14",
		Name:            "Check for prune setting in sync policy",
		PassDescription: "Applications have safe prune settings (prune=false or not specified).",
		FailDescription: "Applications have prune=true in sync policy. This can lead to unintended resource deletion.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule14) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Check if automated sync policy exists
		_, hasAutomated := getYAMLValue(&yamlDoc, "spec.syncPolicy.automated")
		if !hasAutomated {
			// syncPolicy is not automated, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Extract prune value using path-based access
		pruneNode, hasPrune := getYAMLValue(&yamlDoc, "spec.syncPolicy.automated.prune")

		if hasPrune {
			prune, pruneValid := getBoolValue(pruneNode)
			if pruneValid && prune {
				// prune=true - rule fails
				failedFiles = append(failedFiles, app.Name)
			} else {
				// prune=false - rule passes
				passedFiles = append(passedFiles, app.Name)
			}
		} else {
			// prune not mentioned - rule passes
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
