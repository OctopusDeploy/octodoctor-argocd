// Rule 15: Check for hardcoded values in helm directive
// This rule checks if applications have hardcoded values in helm configuration
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
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
			result := checkHelmSource(sourceNode)
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

// helmCheckResult holds the result of checking a single source entry for helm configuration
type helmCheckResult struct {
	hasHelm bool
	passed  bool
}

// checkHelmSource examines a single source entry for helm configuration
func checkHelmSource(sourceNode *yaml.Node) helmCheckResult {
	// Check if helm directive exists
	helmNode, hasHelm := getYAMLValueFromNode(sourceNode, "helm")
	if !hasHelm {
		return helmCheckResult{hasHelm: false, passed: true}
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

	// Rule fails if it has values or valuesObject
	passed := !hasValues && !hasValuesObject
	return helmCheckResult{hasHelm: true, passed: passed}
}
