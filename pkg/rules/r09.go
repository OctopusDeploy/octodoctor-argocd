// Rule 09: Check for duplicate generators in ApplicationSets
// This rule checks if ApplicationSets use matrix/merge generators with duplicate generator types
// which can lead to complex and confusing configurations
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
)

type Rule09 struct {
}

func (r *Rule09) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R09",
		Name:            "Check for duplicate generators in matrix/merge configurations",
		PassDescription: "ApplicationSets should use matrix/merge generators with different generator types.",
		FailDescription: "ApplicationSets have duplicate generators in matrix/merge configurations. This creates complex and confusing configurations that are hard to maintain.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule09) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Get the generators array
		generatorsNode, hasGenerators := getYAMLValue(&yamlDoc, "spec.generators")
		if !hasGenerators || generatorsNode.Kind != yaml.SequenceNode {
			// No generators found, mark as inactive
			inactiveFiles = append(inactiveFiles, appSet.Name)
			continue
		}

		// Check if there's a matrix/merge combination
		hasCombinationGenerator := false
		hasDuplicateGenerators := false

		// Iterate through all generators
		for _, generatorNode := range generatorsNode.Content {
			if generatorNode.Kind != yaml.MappingNode {
				continue
			}

			// Check for matrix or merge generators
			matrixNode, hasMatrix := getYAMLValueFromNode(generatorNode, "matrix")
			mergeNode, hasMerge := getYAMLValueFromNode(generatorNode, "merge")

			if hasMatrix {
				hasCombinationGenerator = true
				if hasDuplicateGeneratorsInCombination(matrixNode) {
					hasDuplicateGenerators = true
				}
			} else if hasMerge {
				hasCombinationGenerator = true
				if hasDuplicateGeneratorsInCombination(mergeNode) {
					hasDuplicateGenerators = true
				}
			}
		}

		// Categorize the ApplicationSet
		if !hasCombinationGenerator {
			// No matrix/merge combination
			inactiveFiles = append(inactiveFiles, appSet.Name)
		} else if hasDuplicateGenerators {
			// Has duplicate generators in matrix/merge
			failedFiles = append(failedFiles, appSet.Name)
		} else {
			// Has matrix/merge with different generators
			passedFiles = append(passedFiles, appSet.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}

// hasDuplicateGeneratorsInCombination checks if a matrix or merge generator has duplicate generator types
func hasDuplicateGeneratorsInCombination(combinationNode *yaml.Node) bool {
	// Get the child generators array
	childGeneratorsNode, hasChildGenerators := getYAMLValueFromNode(combinationNode, "generators")
	if !hasChildGenerators || childGeneratorsNode.Kind != yaml.SequenceNode {
		return false
	}

	// Track generator types we've seen
	generatorTypesSeen := make(map[string]bool)

	// Check each child generator
	for _, childGenNode := range childGeneratorsNode.Content {
		if childGenNode.Kind != yaml.MappingNode {
			continue
		}

		// Get the generator type (the key of the generator)
		generatorType := getGeneratorType(childGenNode)
		if generatorType == "" {
			continue
		}

		// Check if we've seen this generator type before
		if generatorTypesSeen[generatorType] {
			// Found a duplicate
			return true
		}

		// Mark this generator type as seen
		generatorTypesSeen[generatorType] = true
	}

	return false
}

// getGeneratorType extracts the generator type from a generator node
// Returns the generator key (e.g., "git", "list", "clusters", "scmProvider", "pullRequest", "matrix", "merge")
func getGeneratorType(generatorNode *yaml.Node) string {
	if generatorNode.Kind != yaml.MappingNode {
		return ""
	}

	// The generator type is the first key in the mapping
	// Generator nodes have structure like: { "git": {...} } or { "list": {...} }
	if len(generatorNode.Content) >= 2 {
		// Content[0] is the key, Content[1] is the value
		keyNode := generatorNode.Content[0]
		if keyNode.Kind == yaml.ScalarNode {
			return keyNode.Value
		}
	}

	return ""
}
