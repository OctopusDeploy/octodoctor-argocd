// Rule 06: Check if manifest-generate-paths annotation exists
// This rule checks if applications have the argocd.argoproj.io/manifest-generate-paths annotation
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule06 struct {
}

func (r *Rule06) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R06",
		Name:            "Check if manifest-generate-paths annotation exists",
		PassDescription: "Applications have the argocd.argoproj.io/manifest-generate-paths annotation configured.",
		FailDescription: "Applications do not have the argocd.argoproj.io/manifest-generate-paths annotation configured.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule06) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string

	// Iterate over all ArgoCD Application files
	for _, app := range file_analysis.Applications {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(app.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if manifest-generate-paths annotation exists
		annotationsNode, hasAnnotations := getYAMLValue(&yamlDoc, "metadata.annotations")
		hasAnnotation := false
		if hasAnnotations && annotationsNode.Kind == yaml.MappingNode {
			// Look for the specific annotation key in the annotations mapping
			for i := 0; i < len(annotationsNode.Content); i += 2 {
				if annotationsNode.Content[i].Value == "argocd.argoproj.io/manifest-generate-paths" {
					hasAnnotation = true
					break
				}
			}
		}
		if hasAnnotation {
			// manifest-generate-paths annotation exists, rule passes
			passedFiles = append(passedFiles, app.Name)
		} else {
			// manifest-generate-paths annotation doesn't exist, rule fails
			failedFiles = append(failedFiles, app.Name)
		}
	}

	// Determine overall rule result
	// Note: This rule cannot be inactive - it either passes or fails for each application
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, nil)
}
