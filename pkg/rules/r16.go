// Rule 16: Check for hardcoded values in kustomize directive
// This rule checks if applications have kustomize configuration with hardcoded values
package rules

import (
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
)

type Rule16 struct {
}

func (r *Rule16) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R16",
		Name:            "Check for hardcoded values in kustomize directive",
		PassDescription: "Applications do not use hardcoded values in kustomize configuration.",
		FailDescription: "Applications use hardcoded values in kustomize configuration. Consider using kustomize overlays instead of hardcoded values for better maintainability.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule16) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Check all sources for kustomize directive
		hasKustomize := slices.ContainsFunc(sourceNodes, checkKustomizeSource)

		// Categorize the application
		if hasKustomize {
			// Has kustomize directive - rule fails
			failedFiles = append(failedFiles, app.Name)
		} else {
			// No kustomize directive - rule passes
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	// Note: This rule cannot be inactive - it either passes or fails
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}

// checkKustomizeSource examines a single source entry for kustomize configuration
func checkKustomizeSource(sourceNode *yaml.Node) bool {
	// Check if kustomize directive exists
	_, hasKustomize := getYAMLValueFromNode(sourceNode, "kustomize")
	return hasKustomize
}
