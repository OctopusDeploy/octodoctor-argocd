// Rule 01: Check sync policy automated settings
// This rule checks if applications have dangerous sync policy settings
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule01 struct {
}

func (r *Rule01) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R01",
		Name:            "Check sync policy automated settings",
		PassDescription: "Applications have safe sync policy settings (prune=false or allowEmpty=true).",
		FailDescription: "Applications have dangerous sync policy settings (prune=true and allowEmpty=true). This combination can lead to unintended resource deletion.",
		Severity:        SeverityHigh,
	}
}

func (r *Rule01) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Extract prune and allowEmpty values using path-based access
		pruneNode, hasPrune := getYAMLValue(&yamlDoc, "spec.syncPolicy.automated.prune")
		allowEmptyNode, hasAllowEmpty := getYAMLValue(&yamlDoc, "spec.syncPolicy.automated.allowEmpty")

		prune, pruneValid := getBoolValue(pruneNode)
		allowEmpty, allowEmptyValid := getBoolValue(allowEmptyNode)

		// Rule logic:
		// FAIL: if prune=true AND allowEmpty=true
		// PASS: if prune=false OR allowEmpty=false (or missing values)
		if hasPrune && hasAllowEmpty && pruneValid && allowEmptyValid && prune && allowEmpty {
			// Dangerous combination: prune=true and allowEmpty=true
			failedFiles = append(failedFiles, app.Name)
		} else {
			// Safe combination: either prune=false or allowEmpty=true (or missing values)
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
