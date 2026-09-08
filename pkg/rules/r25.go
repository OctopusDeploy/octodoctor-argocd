// Rule 25: Check for applications with more than 2 sources
// This rule checks if multi-source applications have more than 2 sources
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule25 struct {
}

func (r *Rule25) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R25",
		Name:            "Check for applications with more than 2 sources",
		DocSlug:         "r25-more-than-two-sources",
		PassDescription: "Applications have at most 2 sources in multi-source configuration.",
		FailDescription: "Applications with more than 2 sources need a different grouping mechanism.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule25) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Check if application has single source (spec.source)
		_, hasSingleSource := getYAMLValue(&yamlDoc, "spec.source")
		if hasSingleSource {
			// Single source application - rule is inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Check if application has multi-source (spec.sources)
		sourcesNode, hasMultiSource := getYAMLValue(&yamlDoc, "spec.sources")
		if !hasMultiSource {
			// No sources found - mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Check if sources is a sequence (array)
		if sourcesNode.Kind != yaml.SequenceNode {
			// sources exists but is not an array, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Count the number of sources
		numSources := len(sourcesNode.Content)

		// Rule logic:
		// FAIL: if sources array has more than 2 items
		// PASS: if sources array has at most 2 items
		if numSources > 2 {
			// More than 2 sources - rule fails
			failedFiles = append(failedFiles, app.Name)
		} else {
			// At most 2 sources - rule passes
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
