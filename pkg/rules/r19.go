// Rule 19: Check targetRevision field usage
// This rule checks if applications use appropriate targetRevision values
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule19 struct {
}

func (r *Rule19) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R19",
		Name:            "Check targetRevision field usage",
		PassDescription: "Applications have targetRevision set to HEAD or not specified.",
		FailDescription: "Applications have targetRevision set to a specific value (not HEAD) for Git repositories. This can lead to stale deployments.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule19) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
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
		loadsHelmChart := false
		rulePassed := true
		for _, sourceNode := range sourceNodes {
			result := checkSource(sourceNode)
			if result.loadsHelmChart {
				loadsHelmChart = true
				break
			}
			if !result.passed {
				rulePassed = false
			}
		}

		// Categorize the application
		if loadsHelmChart {
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

// sourceCheckResult holds the result of checking a single source entry
type sourceCheckResult struct {
	loadsHelmChart bool
	passed         bool
}

// checkSource examines a single source entry and returns whether it loads a Helm chart from a repository
func checkSource(sourceNode *yaml.Node) sourceCheckResult {
	// Check if chart property exists
	_, loadsHelmChart := getYAMLValueFromNode(sourceNode, "chart")
	if loadsHelmChart {
		return sourceCheckResult{loadsHelmChart: true, passed: true}
	}

	// Check targetRevision
	targetRevisionNode, hasTargetRevision := getYAMLValueFromNode(sourceNode, "targetRevision")
	if !hasTargetRevision {
		// targetRevision doesn't exist, passes by default
		return sourceCheckResult{loadsHelmChart: false, passed: true}
	}

	// Check if targetRevision is HEAD
	targetRevision := targetRevisionNode.Value
	passed := targetRevision == "HEAD"
	return sourceCheckResult{loadsHelmChart: false, passed: passed}
}
