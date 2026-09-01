// Rule 20: Check if ApplicationSets use cluster generators with multiple clusters
// This rule checks if ApplicationSets properly use the cluster generator when multiple clusters are configured
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
)

type Rule20 struct {
}

func (r *Rule20) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R20",
		Name:            "Check cluster generator usage with multiple clusters",
		PassDescription: "ApplicationSets correctly use cluster generators when multiple clusters are available.",
		FailDescription: "ApplicationSets do not use cluster generators despite having multiple connected clusters. This can lead to performance issues and configuration complexity.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule20) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	// Step 1: Count the number of external clusters by checking secrets with cluster label
	externalClusterCount := countExternalClusters(file_analysis)

	// If there are no external clusters, the rule is inactive
	if externalClusterCount == 0 {
		result := analysis.RuleExecution{
			ID:             r.Metadata().ID,
			Name:           r.Metadata().Name,
			Description:    r.Metadata().PassDescription,
			Severity:       r.Metadata().Severity,
			FilesTriggered: nil,
		}
		return RuleResult{
			Active:    false,
			Passed:    true,
			Execution: result,
		}, nil
	}

	// If there is only 1 external cluster, the rule passes automatically
	if externalClusterCount == 1 {
		result := analysis.RuleExecution{
			ID:             r.Metadata().ID,
			Name:           r.Metadata().Name,
			Description:    r.Metadata().PassDescription,
			Severity:       r.Metadata().Severity,
			FilesTriggered: nil,
		}
		return RuleResult{
			Active:    true,
			Passed:    true,
			Execution: result,
		}, nil
	}

	// Step 2: If there are 2+ clusters, check if ApplicationSets use cluster generators
	var failedFiles []string
	var passedFiles []string

	for _, appSet := range file_analysis.ApplicationSets {
		// Parse YAML content
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(appSet.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if the ApplicationSet has a cluster generator
		hasClusterGenerator := hasClusterGeneratorInAppSet(&yamlDoc)

		if hasClusterGenerator {
			passedFiles = append(passedFiles, appSet.Name)
		} else {
			failedFiles = append(failedFiles, appSet.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, []string{})
}

// countExternalClusters counts the number of secrets with the cluster-type label
// as these represent external clusters in Argo CD
func countExternalClusters(file_analysis *analysis.ArgoCDAnalysis) int {
	clusterCount := 0

	for _, secret := range file_analysis.Secrets {
		// Parse YAML content
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(secret.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check for the label argocd.argoproj.io/secret-type: cluster
		labelsNode, hasLabels := getYAMLValue(&yamlDoc, "metadata.labels")
		if !hasLabels || labelsNode.Kind != yaml.MappingNode {
			continue
		}

		// Look for the secret-type label
		for i := 0; i < len(labelsNode.Content); i += 2 {
			if labelsNode.Content[i].Value == "argocd.argoproj.io/secret-type" {
				valueNode := labelsNode.Content[i+1]
				if valueNode.Kind == yaml.ScalarNode && valueNode.Value == "cluster" {
					clusterCount++
					break
				}
			}
		}
	}

	return clusterCount
}

// hasClusterGeneratorInAppSet checks if an ApplicationSet has a cluster generator at any level
func hasClusterGeneratorInAppSet(yamlDoc *yaml.Node) bool {
	// Get the generators array
	generatorsNode, hasGenerators := getYAMLValue(yamlDoc, "spec.generators")
	if !hasGenerators || generatorsNode.Kind != yaml.SequenceNode {
		return false
	}

	// Check if any generator is a cluster generator or contains one
	return hasClusterGeneratorInList(generatorsNode)
}

// hasClusterGeneratorInList recursively checks if a generators list contains a cluster generator
func hasClusterGeneratorInList(generatorsNode *yaml.Node) bool {
	for _, generatorNode := range generatorsNode.Content {
		if generatorNode.Kind != yaml.MappingNode {
			continue
		}

		// Check if this is a direct cluster generator
		if hasClusterGenerator(generatorNode) {
			return true
		}

		// Check within matrix generators
		matrixNode, hasMatrix := getYAMLValueFromNode(generatorNode, "matrix")
		if hasMatrix && matrixNode.Kind == yaml.MappingNode {
			childGeneratorsNode, hasChildGenerators := getYAMLValueFromNode(matrixNode, "generators")
			if hasChildGenerators && childGeneratorsNode.Kind == yaml.SequenceNode {
				if hasClusterGeneratorInList(childGeneratorsNode) {
					return true
				}
			}
		}

		// Check within merge generators
		mergeNode, hasMerge := getYAMLValueFromNode(generatorNode, "merge")
		if hasMerge && mergeNode.Kind == yaml.MappingNode {
			childGeneratorsNode, hasChildGenerators := getYAMLValueFromNode(mergeNode, "generators")
			if hasChildGenerators && childGeneratorsNode.Kind == yaml.SequenceNode {
				if hasClusterGeneratorInList(childGeneratorsNode) {
					return true
				}
			}
		}
	}

	return false
}

// hasClusterGenerator checks if a generator node is a cluster generator
func hasClusterGenerator(generatorNode *yaml.Node) bool {
	_, hasClusters := getYAMLValueFromNode(generatorNode, "clusters")
	return hasClusters
}
