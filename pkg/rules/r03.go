// Rule 03: Check for requeueAfterSeconds setting
// This rule checks if ApplicationSets have requeueAfterSeconds set to less than 3 minutes (180 seconds)
// which can cause excessive API calls and rate limiting issues
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule03 struct {
}

const minRequeueAfterSeconds = 180 // 3 minutes

func (r *Rule03) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R03",
		Name:            "Check requeueAfterSeconds setting",
		DocSlug:         "r03-low-requeue-after-seconds",
		PassDescription: "ApplicationSets have requeueAfterSeconds set to 180 seconds (3 minutes) or more.",
		FailDescription: "ApplicationSets have requeueAfterSeconds set to less than 180 seconds (3 minutes). This can cause excessive API calls and rate limiting issues, especially for GitHub.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule03) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	// Iterate over all ArgoCD ApplicationSet files
	for _, appSet := range file_analysis.ApplicationSets {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(appSet.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Recursively search for all requeueAfterSeconds values in the entire YAML
		requeueValues := findAllRequeueAfterSeconds(&yamlDoc)

		// Categorize the ApplicationSet
		if len(requeueValues) == 0 {
			// No requeueAfterSeconds found anywhere
			inactiveFiles = append(inactiveFiles, appSet.Name)
		} else {
			// Check if any value is below the minimum threshold
			hasLowValue := false
			for _, value := range requeueValues {
				if value < minRequeueAfterSeconds {
					hasLowValue = true
					break
				}
			}

			if hasLowValue {
				// At least one requeueAfterSeconds is too low
				failedFiles = append(failedFiles, appSet.Name)
			} else {
				// All requeueAfterSeconds values are acceptable (>= 180 seconds)
				passedFiles = append(passedFiles, appSet.Name)
			}
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}

// findAllRequeueAfterSeconds recursively searches the entire YAML tree for all occurrences of requeueAfterSeconds
// Returns a slice of all found values
func findAllRequeueAfterSeconds(node *yaml.Node) []int {
	var values []int

	if node == nil {
		return values
	}

	switch node.Kind {
	case yaml.DocumentNode:
		// Document node - recurse into content
		for _, child := range node.Content {
			values = append(values, findAllRequeueAfterSeconds(child)...)
		}

	case yaml.MappingNode:
		// Mapping node - check key-value pairs
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			// Check if this key is "requeueAfterSeconds"
			if keyNode.Kind == yaml.ScalarNode && keyNode.Value == "requeueAfterSeconds" {
				// Try to decode the value as an integer
				if valueNode.Kind == yaml.ScalarNode {
					var value int
					if err := valueNode.Decode(&value); err == nil {
						values = append(values, value)
					}
				}
			}

			// Recurse into the value to find nested requeueAfterSeconds
			values = append(values, findAllRequeueAfterSeconds(valueNode)...)
		}

	case yaml.SequenceNode:
		// Sequence node - recurse into each element
		for _, child := range node.Content {
			values = append(values, findAllRequeueAfterSeconds(child)...)
		}
	}

	return values
}
