// Rule 04: Check sync policy replace setting
// This rule checks if applications have replace=true in syncOptions
package rules

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule04 struct {
}

func (r *Rule04) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R04",
		Name:            "Check sync policy replace setting",
		DocSlug:         "r04-sync-policy-replace",
		PassDescription: "Applications have safe sync policy settings (replace=false or replace not specified).",
		FailDescription: "Applications have sync policy replace=true. This can lead to unintended resource deletion.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule04) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Check if syncOptions exists
		syncOptionsNode, hasSyncOptions := getYAMLValue(&yamlDoc, "spec.syncPolicy.syncOptions")
		if !hasSyncOptions {
			// No syncOptions, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Check if syncOptions is a sequence (array)
		if syncOptionsNode.Kind != yaml.SequenceNode {
			// syncOptions exists but is not an array, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Look for Replace setting in syncOptions
		foundReplace := false
		replaceValue := false

		for _, optionNode := range syncOptionsNode.Content {
			if optionNode.Kind == yaml.ScalarNode {
				option := optionNode.Value
				// Check if this option is a Replace setting
				if strings.HasPrefix(option, "Replace=") {
					foundReplace = true
					replaceValueStr := strings.TrimPrefix(option, "Replace=")
					replaceValue = strings.ToLower(replaceValueStr) == "true"
					break
				}
			}
		}

		if !foundReplace {
			// Replace not mentioned in syncOptions, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Rule logic:
		// FAIL: if Replace=true
		// PASS: if Replace=false
		if replaceValue {
			// Replace=true is dangerous
			failedFiles = append(failedFiles, app.Name)
		} else {
			// Replace=false is safe
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
