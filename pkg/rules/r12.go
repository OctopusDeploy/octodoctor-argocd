// Rule 12: Check for excessive nesting of generators in ApplicationSets
// This rule checks if ApplicationSets have 3 or more levels of nested generators
// which creates overly complex configurations
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule12 struct {
}

func (r *Rule12) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R12",
		Name:            "Check for excessive generator nesting levels",
		PassDescription: "ApplicationSets use 2 levels of generators (single matrix or merge combination).",
		FailDescription: "ApplicationSets have 3 or more levels of generators by nesting matrix/merge combinations. This creates overly complex and hard-to-maintain configurations.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule12) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Check if there's a matrix/merge combination and its nesting level
		hasCombinationGenerator := false
		hasNestedCombination := false

		// Iterate through all top-level generators
		for _, generatorNode := range generatorsNode.Content {
			if generatorNode.Kind != yaml.MappingNode {
				continue
			}

			// Check for matrix or merge generators
			matrixNode, hasMatrix := getYAMLValueFromNode(generatorNode, "matrix")
			mergeNode, hasMerge := getYAMLValueFromNode(generatorNode, "merge")

			if hasMatrix {
				hasCombinationGenerator = true
				if hasNestedCombinationGenerator(matrixNode) {
					hasNestedCombination = true
				}
			} else if hasMerge {
				hasCombinationGenerator = true
				if hasNestedCombinationGenerator(mergeNode) {
					hasNestedCombination = true
				}
			}
		}

		// Categorize the ApplicationSet
		if !hasCombinationGenerator {
			// No matrix/merge combination at all
			inactiveFiles = append(inactiveFiles, appSet.Name)
		} else if hasNestedCombination {
			// Has nested matrix/merge (3+ levels)
			failedFiles = append(failedFiles, appSet.Name)
		} else {
			// Has matrix/merge but no nesting (2 levels)
			passedFiles = append(passedFiles, appSet.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}

// hasNestedCombinationGenerator checks if any child generator is also a matrix or merge generator
// This indicates 3 or more levels of nesting
func hasNestedCombinationGenerator(combinationNode *yaml.Node) bool {
	// Get the child generators array
	childGeneratorsNode, hasChildGenerators := getYAMLValueFromNode(combinationNode, "generators")
	if !hasChildGenerators || childGeneratorsNode.Kind != yaml.SequenceNode {
		return false
	}

	// Check each child generator to see if it's also a matrix or merge
	for _, childGenNode := range childGeneratorsNode.Content {
		if childGenNode.Kind != yaml.MappingNode {
			continue
		}

		// Check if this child is a matrix or merge generator
		_, hasMatrix := getYAMLValueFromNode(childGenNode, "matrix")
		_, hasMerge := getYAMLValueFromNode(childGenNode, "merge")

		if hasMatrix || hasMerge {
			// Found a nested combination generator (3+ levels)
			return true
		}
	}

	return false
}
