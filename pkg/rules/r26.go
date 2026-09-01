// Rule 26: Check for applications with multiple sources containing Helm charts
// This rule checks if multi-source applications have 2 or more sources with chart property
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
)

type Rule26 struct {
}

func (r *Rule26) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R26",
		Name:            "Check for applications with multiple sources containing Helm charts",
		PassDescription: "Applications have at most one source with a Helm chart in multi-source configuration.",
		FailDescription: "Applications have 2 or more sources with Helm charts. Multiple Helm charts should be grouped differently.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule26) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Count how many sources have the chart property
		chartCount := 0
		for _, sourceNode := range sourcesNode.Content {
			_, hasChart := getYAMLValueFromNode(sourceNode, "chart")
			if hasChart {
				chartCount++
			}
		}

		// Rule logic:
		// FAIL: if 2 or more sources have chart property
		// PASS: if exactly 1 source has chart property (or 0 charts)
		if chartCount >= 2 {
			// 2 or more charts - rule fails
			failedFiles = append(failedFiles, app.Name)
		} else {
			// At most 1 chart - rule passes
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
