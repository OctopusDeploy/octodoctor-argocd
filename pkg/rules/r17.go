// Rule 17: Check for parameter overrides in helm directive
// This rule checks if applications have parameter overrides in helm configuration
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule17 struct {
}

func (r *Rule17) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R17",
		Name:            "Check for parameter overrides in helm directive",
		DocSlug:         "r17-parameter-overrides",
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

		// Get source nodes (either single source or multi-source)
		var sourceNodes []*yaml.Node
		sourcesNode, hasMultiSource := getYAMLValue(&yamlDoc, "spec.sources")
		if hasMultiSource && sourcesNode.Kind == yaml.SequenceNode {
			// Multi-source application
			sourceNodes = sourcesNode.Content
		} else {
			// Single source application
			if sourceNode, hasSingleSource := getYAMLValue(&yamlDoc, "spec.source"); hasSingleSource {
				sourceNodes = []*yaml.Node{sourceNode}
			}
		}

		// Check all sources
		hasHelm := false
		rulePassed := true
		for _, sourceNode := range sourceNodes {
			result := checkHelmParametersSource(sourceNode)
			if result.hasHelm {
				hasHelm = true
				if !result.passed {
					rulePassed = false
				}
			}
		}

		// Categorize the application
		if !hasHelm {
			// No helm directive in any source
			inactiveFiles = append(inactiveFiles, app.Name)
		} else if rulePassed {
			passedFiles = append(passedFiles, app.Name)
		} else {
			failedFiles = append(failedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}

// helmParametersCheckResult holds the result of checking a single source entry for helm parameters
type helmParametersCheckResult struct {
	hasHelm bool
	passed  bool
}

// checkHelmParametersSource examines a single source entry for helm parameter configuration
func checkHelmParametersSource(sourceNode *yaml.Node) helmParametersCheckResult {
	// Check if helm directive exists
	_, hasHelm := getYAMLValueFromNode(sourceNode, "helm")
	if !hasHelm {
		return helmParametersCheckResult{hasHelm: false, passed: true}
	}

	// Check if helm directive contains parameters field
	_, hasParameters := getYAMLValueFromNode(sourceNode, "helm.parameters")

	// Rule fails if it has parameters
	passed := !hasParameters
	return helmParametersCheckResult{hasHelm: true, passed: passed}
}
