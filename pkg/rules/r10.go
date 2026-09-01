// Rule 10: Check argocd-cmd-params-cm sync timeout configuration
// This rule checks if the controller sync timeout is properly configured (not 0 seconds)
package rules

import (
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
)

type Rule10 struct {
}

func (r *Rule10) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R10",
		Name:            "Check argocd-cmd-params-cm sync timeout configuration",
		PassDescription: "ArgoCD controller sync timeout is properly configured (not 0 seconds).",
		FailDescription: "ArgoCD controller sync timeout is set to 0 seconds, which means no timeout and can cause inefficient synchronization.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule10) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	foundArgoCdCmdParamsCm := false

	// Iterate over all ArgoCD ConfigMap files
	for _, configMap := range file_analysis.ConfigMaps {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(configMap.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if this is the argocd-cmd-params-cm ConfigMap
		nameNode, hasName := getYAMLValue(&yamlDoc, "metadata.name")
		if !hasName || nameNode.Kind != yaml.ScalarNode || nameNode.Value != "argocd-cmd-params-cm" {
			// Not the argocd-cmd-params-cm ConfigMap, skip this file
			inactiveFiles = append(inactiveFiles, configMap.Name)
			continue
		}

		foundArgoCdCmdParamsCm = true

		// Check if controller.sync.timeout.seconds exists in data
		// Note: controller.sync.timeout.seconds is actually a single key with dots, not nested
		dataNode, hasData := getYAMLValue(&yamlDoc, "data")
		var timeoutNode *yaml.Node
		hasTimeout := false

		if hasData && dataNode.Kind == yaml.MappingNode {
			// Look for the controller.sync.timeout.seconds key directly in the data mapping
			for i := 0; i < len(dataNode.Content); i += 2 {
				keyNode := dataNode.Content[i]
				valueNode := dataNode.Content[i+1]
				if keyNode.Value == "controller.sync.timeout.seconds" {
					timeoutNode = valueNode
					hasTimeout = true
					break
				}
			}
		}

		if !hasTimeout {
			// Property doesn't exist, consider as pass (no explicit timeout set)
			passedFiles = append(passedFiles, configMap.Name)
			continue
		}

		// Check if the timeout value is 0
		if timeoutNode.Kind == yaml.ScalarNode {
			timeoutValue := timeoutNode.Value
			// Try to parse as integer
			if timeoutInt, err := strconv.Atoi(timeoutValue); err == nil {
				if timeoutInt == 0 {
					// Timeout is 0, rule fails
					failedFiles = append(failedFiles, configMap.Name)
				} else {
					// Timeout is not 0, rule passes
					passedFiles = append(passedFiles, configMap.Name)
				}
			} else {
				// Can't parse as integer, consider as pass (non-numeric value)
				passedFiles = append(passedFiles, configMap.Name)
			}
		} else {
			// Not a scalar value, consider as pass
			passedFiles = append(passedFiles, configMap.Name)
		}
	}

	// If no argocd-cmd-params-cm ConfigMap was found, mark the rule as inactive
	if !foundArgoCdCmdParamsCm {
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
