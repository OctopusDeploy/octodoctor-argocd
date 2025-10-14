// Rule 18: Check for ApplicationSets combining SCM and Pull Request generators
// This rule checks if ApplicationSets use matrix/merge generators that combine SCM Provider and Pull Request generators
// which can cause excessive API calls and rate limiting issues
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule18 struct {
}

func (r *Rule18) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R18",
		Name:            "Check for ApplicationSets combining SCM and Pull Request generators",
		PassDescription: "ApplicationSets use Pull Request generators safely without combining with SCM Provider generators.",
		FailDescription: "ApplicationSets combine SCM Provider and Pull Request generators in matrix/merge configurations. This can cause exponential API calls and exhaust rate limits.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule18) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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

		// Check if there's a Pull Request generator in a matrix/merge combination
		hasPRInCombination := false
		hasDangerousCombination := false

		// Iterate through all generators
		for _, generatorNode := range generatorsNode.Content {
			if generatorNode.Kind != yaml.MappingNode {
				continue
			}

			// Check for matrix or merge generators
			matrixNode, hasMatrix := getYAMLValueFromNode(generatorNode, "matrix")
			mergeNode, hasMerge := getYAMLValueFromNode(generatorNode, "merge")

			if hasMatrix {
				result := checkCombinationGenerator(matrixNode)
				hasPRInCombination = hasPRInCombination || result.hasPRGenerator
				hasDangerousCombination = hasDangerousCombination || result.isDangerousGeneratorCombination
			} else if hasMerge {
				result := checkCombinationGenerator(mergeNode)
				hasPRInCombination = hasPRInCombination || result.hasPRGenerator
				hasDangerousCombination = hasDangerousCombination || result.isDangerousGeneratorCombination
			}
		}

		// Categorize the ApplicationSet
		if !hasPRInCombination {
			// No PR generator in matrix/merge combination (either no PR at all, or PR is standalone)
			inactiveFiles = append(inactiveFiles, appSet.Name)
		} else if hasDangerousCombination {
			// Has PR generator combined with SCM in matrix/merge
			failedFiles = append(failedFiles, appSet.Name)
		} else {
			// Has PR generator in matrix/merge but not combined with SCM
			passedFiles = append(passedFiles, appSet.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}

// generatorCheckResult holds the result of checking a combination generator
type generatorCheckResult struct {
	hasPRGenerator                  bool // Whether it contains a PR generator
	isDangerousGeneratorCombination bool // Whether it's the dangerous SCM + PR combination
}

// checkCombinationGenerator examines a matrix or merge generator for dangerous combinations
func checkCombinationGenerator(combinationNode *yaml.Node) generatorCheckResult {
	result := generatorCheckResult{hasPRGenerator: false, isDangerousGeneratorCombination: false}

	// Get the child generators array
	childGeneratorsNode, hasChildGenerators := getYAMLValueFromNode(combinationNode, "generators")
	if !hasChildGenerators || childGeneratorsNode.Kind != yaml.SequenceNode {
		return result
	}

	// Check what types of generators are present in the combination
	hasSCM := false
	hasPR := false

	for _, childGenNode := range childGeneratorsNode.Content {
		if childGenNode.Kind != yaml.MappingNode {
			continue
		}

		// Check for SCM Provider generator
		if hasSCMProviderGenerator(childGenNode) {
			hasSCM = true
		}

		// Check for Pull Request generator
		if hasPullRequestGenerator(childGenNode) {
			hasPR = true
		}
	}

	result.hasPRGenerator = hasPR
	result.isDangerousGeneratorCombination = hasSCM && hasPR

	return result
}

// hasSCMProviderGenerator checks if a generator node contains an SCM Provider generator
func hasSCMProviderGenerator(generatorNode *yaml.Node) bool {
	_, hasSCMProvider := getYAMLValueFromNode(generatorNode, "scmProvider")
	return hasSCMProvider
}

// hasPullRequestGenerator checks if a generator node contains a Pull Request generator
func hasPullRequestGenerator(generatorNode *yaml.Node) bool {
	_, hasPR := getYAMLValueFromNode(generatorNode, "pullRequest")
	return hasPR
}
