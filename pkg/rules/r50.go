// Rule 50: Check argocd-cm resourceTrackingMethod configuration
// This rule checks if resourceTrackingMethod is set to one of the accepted values (label, annotation, or annotation+label)
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
)

type Rule50 struct {
}

func (r *Rule50) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R50",
		Name:            "Check argocd-cm resourceTrackingMethod configuration",
		PassDescription: "ArgoCD resourceTrackingMethod is set to one of the accepted values (label, annotation, or annotation+label).",
		FailDescription: "ArgoCD resourceTrackingMethod is set to an invalid value. It must be exactly one of: label, annotation, or annotation+label.",
		Severity:        SeverityHigh,
	}
}

func (r *Rule50) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	foundArgoCdCm := false

	// Accepted values for resourceTrackingMethod
	acceptedValues := map[string]bool{
		"label":            true,
		"annotation":       true,
		"annotation+label": true,
	}

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

		// Check if application.resourceTrackingMethod exists in data
		// Note: application.resourceTrackingMethod is actually a single key with dots, not nested
		dataNode, hasData := getYAMLValue(&yamlDoc, "data")
		var trackingMethodNode *yaml.Node
		hasTrackingMethod := false

		if hasData && dataNode.Kind == yaml.MappingNode {
			// Look for the application.resourceTrackingMethod key directly in the data mapping
			for i := 0; i < len(dataNode.Content); i += 2 {
				keyNode := dataNode.Content[i]
				valueNode := dataNode.Content[i+1]
				if keyNode.Value == "application.resourceTrackingMethod" {
					trackingMethodNode = valueNode
					hasTrackingMethod = true
					break
				}
			}
		}

		if !hasTrackingMethod {
			// Property doesn't exist, consider as inactive (rule doesn't apply)
			inactiveFiles = append(inactiveFiles, configMap.Name)
			continue
		}

		// Check if the tracking method value is one of the accepted values
		if trackingMethodNode.Kind == yaml.ScalarNode {
			trackingMethodValue := trackingMethodNode.Value
			// Check if the value matches exactly one of the accepted values
			if acceptedValues[trackingMethodValue] {
				// Value is accepted, rule passes
				passedFiles = append(passedFiles, configMap.Name)
			} else {
				// Value is not accepted, rule fails
				failedFiles = append(failedFiles, configMap.Name)
			}
		} else {
			// Not a scalar value, consider as fail (invalid format)
			failedFiles = append(failedFiles, configMap.Name)
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
