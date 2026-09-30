// Rule 11: Check argocd-cm anonymous user access configuration
// This rule checks if anonymous user access is disabled for security
package rules

import (
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule11 struct {
}

func (r *Rule11) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R11",
		Name:            "Check argocd-cm anonymous user access configuration",
		DocSlug:         "r11-anonymous-access",
		PassDescription: "Argo CD anonymous user access is disabled for security.",
		FailDescription: "Argo CD anonymous user access is enabled, which poses a security risk as it lets unauthenticated users access the system.",
		Severity:        SeverityHigh,
	}
}

func (r *Rule11) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	foundArgoCdCm := false

	// Iterate over all ArgoCD ConfigMap files
	for _, configMap := range file_analysis.ConfigMaps {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(configMap.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if this is the argocd-cm ConfigMap
		nameNode, hasName := getYAMLValue(&yamlDoc, "metadata.name")
		if !hasName || nameNode.Kind != yaml.ScalarNode || nameNode.Value != "argocd-cm" {
			// Not the argocd-cm ConfigMap, skip this file
			inactiveFiles = append(inactiveFiles, configMap.Name)
			continue
		}

		foundArgoCdCm = true

		// Check if users.anonymous.enabled exists in data
		// Note: users.anonymous.enabled is actually a single key with dots, not nested
		dataNode, hasData := getYAMLValue(&yamlDoc, "data")
		var anonymousNode *yaml.Node
		hasAnonymous := false

		if hasData && dataNode.Kind == yaml.MappingNode {
			// Look for the users.anonymous.enabled key directly in the data mapping
			for i := 0; i < len(dataNode.Content); i += 2 {
				keyNode := dataNode.Content[i]
				valueNode := dataNode.Content[i+1]
				if keyNode.Value == "users.anonymous.enabled" {
					anonymousNode = valueNode
					hasAnonymous = true
					break
				}
			}
		}

		if !hasAnonymous {
			// Property doesn't exist, consider as pass (anonymous access disabled by default)
			passedFiles = append(passedFiles, configMap.Name)
			continue
		}

		// Check if the anonymous value is true
		if anonymousNode.Kind == yaml.ScalarNode {
			anonymousValue := anonymousNode.Value
			// Try to parse as boolean
			if anonymousBool, err := strconv.ParseBool(anonymousValue); err == nil {
				if anonymousBool {
					// Anonymous access is enabled, rule fails
					failedFiles = append(failedFiles, configMap.Name)
				} else {
					// Anonymous access is disabled, rule passes
					passedFiles = append(passedFiles, configMap.Name)
				}
			} else {
				// Can't parse as boolean, consider as pass (non-boolean value likely means disabled)
				passedFiles = append(passedFiles, configMap.Name)
			}
		} else {
			// Not a scalar value, consider as pass
			passedFiles = append(passedFiles, configMap.Name)
		}
	}

	// If no argocd-cm ConfigMap was found, mark the rule as inactive
	if !foundArgoCdCm {
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

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
