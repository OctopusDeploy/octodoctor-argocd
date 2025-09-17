// Rule 05: Check AppProject sourceRepos settings
// This rule checks if AppProjects have overly permissive sourceRepos settings
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule05 struct {
}

func (r *Rule05) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R05",
		Name:            "Check AppProject sourceRepos settings",
		PassDescription: "AppProjects have restricted sourceRepos settings (not using wildcard '*').",
		FailDescription: "AppProjects have overly permissive sourceRepos settings (using wildcard '*'). This allows applications to deploy from any Git repository, which poses a security risk.",
		Severity:        "high",
	}
}

func (r *Rule05) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	// Iterate over all ArgoCD AppProject files
	for _, project := range file_analysis.AppProjects {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(project.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if sourceRepos exists
		sourceReposNode, hasSourceRepos := getYAMLValue(&yamlDoc, "spec.sourceRepos")
		if !hasSourceRepos {
			// sourceRepos is not present, mark as inactive
			inactiveFiles = append(inactiveFiles, project.Name)
			continue
		}

		// Check if sourceRepos is a sequence (array)
		if sourceReposNode.Kind != yaml.SequenceNode {
			// sourceRepos is not an array, mark as inactive
			inactiveFiles = append(inactiveFiles, project.Name)
			continue
		}

		// Check each item in the sourceRepos array
		hasWildcard := false
		for _, item := range sourceReposNode.Content {
			if item.Kind == yaml.ScalarNode && item.Value == "*" {
				hasWildcard = true
				break
			}
		}

		// Rule logic:
		// FAIL: if any sourceRepo contains "*"
		// PASS: if sourceRepos exists and contains no "*"
		if hasWildcard {
			failedFiles = append(failedFiles, project.Name)
		} else {
			passedFiles = append(passedFiles, project.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
