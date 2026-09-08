// Rule 07: Check ARGOCD_RECONCILIATION_TIMEOUT=0 usage with ARGOCD_DEFAULT_CACHE_EXPIRATION
// This rule checks if reconciliation timeout is set to 0 without proper cache expiration settings
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule07 struct {
}

func (r *Rule07) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R07",
		Name:            "Check ARGOCD_RECONCILIATION_TIMEOUT=0 usage with ARGOCD_DEFAULT_CACHE_EXPIRATION ",
		DocSlug:         "r07-reconciliation-timeout-without-cache-expiration",
		PassDescription: "ArgoCD reconciliation timeout is properly configured or has appropriate cache expiration settings.",
		FailDescription: "Using ARGOCD_RECONCILIATION_TIMEOUT=0 without proper ARGOCD_DEFAULT_CACHE_EXPIRATION  settings can cause performance issues.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule07) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	// Step 1: Find timeout.reconciliation value from argocd-cm ConfigMap
	reconciliationTimeout, foundReconciliation := r.getReconciliationTimeout(file_analysis)

	// If timeout.reconciliation is not present at all, rule is inactive
	if !foundReconciliation {
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

	// If timeout.reconciliation is present and is NOT "0", rule passes
	if !isZeroValue(reconciliationTimeout) {
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

	// Step 2: timeout.reconciliation is "0", check cache expiration
	cacheExpiration, foundCache := r.getCacheExpiration(file_analysis)

	// If cache expiration is present and NOT the default "24h0m0s", rule passes
	if foundCache && cacheExpiration != "24h0m0s" {
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

	// Rule fails: timeout.reconciliation is "0" and cache expiration is missing or default
	var failedFiles []string
	for _, configMap := range file_analysis.ConfigMaps {
		var yamlDoc yaml.Node
		if err := yaml.Unmarshal([]byte(configMap.Content), &yamlDoc); err != nil {
			continue
		}
		nameNode, hasName := getYAMLValue(&yamlDoc, "metadata.name")
		if hasName && nameNode.Kind == yaml.ScalarNode {
			if nameNode.Value == "argocd-cm" || nameNode.Value == "argocd-cmd-params-cm" {
				failedFiles = append(failedFiles, configMap.Name)
			}
		}
	}

	result := analysis.RuleExecution{
		ID:             r.Metadata().ID,
		Name:           r.Metadata().Name,
		Description:    r.Metadata().FailDescription,
		Severity:       r.Metadata().Severity,
		FilesTriggered: failedFiles,
	}
	return RuleResult{
		Active:    true,
		Passed:    false,
		Execution: result,
	}, nil
}

// getReconciliationTimeout extracts timeout.reconciliation from argocd-cm ConfigMap
func (r *Rule07) getReconciliationTimeout(file_analysis *analysis.ArgoCDAnalysis) (string, bool) {
	for _, configMap := range file_analysis.ConfigMaps {
		var yamlDoc yaml.Node
		if err := yaml.Unmarshal([]byte(configMap.Content), &yamlDoc); err != nil {
			continue
		}

		// Check if this is argocd-cm
		nameNode, hasName := getYAMLValue(&yamlDoc, "metadata.name")
		if !hasName || nameNode.Kind != yaml.ScalarNode || nameNode.Value != "argocd-cm" {
			continue
		}

		// Look for timeout.reconciliation in data
		dataNode, hasData := getYAMLValue(&yamlDoc, "data")
		if !hasData || dataNode.Kind != yaml.MappingNode {
			continue
		}

		// Search for timeout.reconciliation key
		for i := 0; i < len(dataNode.Content); i += 2 {
			keyNode := dataNode.Content[i]
			if keyNode.Value == "timeout.reconciliation" {
				valueNode := dataNode.Content[i+1]
				if valueNode.Kind == yaml.ScalarNode {
					return valueNode.Value, true
				}
			}
		}
	}
	return "", false
}

// getCacheExpiration extracts reposerver.default.cache.expiration from argocd-cmd-params-cm ConfigMap
func (r *Rule07) getCacheExpiration(file_analysis *analysis.ArgoCDAnalysis) (string, bool) {
	for _, configMap := range file_analysis.ConfigMaps {
		var yamlDoc yaml.Node
		if err := yaml.Unmarshal([]byte(configMap.Content), &yamlDoc); err != nil {
			continue
		}

		// Check if this is argocd-cmd-params-cm
		nameNode, hasName := getYAMLValue(&yamlDoc, "metadata.name")
		if !hasName || nameNode.Kind != yaml.ScalarNode || nameNode.Value != "argocd-cmd-params-cm" {
			continue
		}

		// Look for reposerver.default.cache.expiration in data
		dataNode, hasData := getYAMLValue(&yamlDoc, "data")
		if !hasData || dataNode.Kind != yaml.MappingNode {
			continue
		}

		// Search for reposerver.default.cache.expiration key
		for i := 0; i < len(dataNode.Content); i += 2 {
			keyNode := dataNode.Content[i]
			if keyNode.Value == "reposerver.default.cache.expiration" {
				valueNode := dataNode.Content[i+1]
				if valueNode.Kind == yaml.ScalarNode {
					return valueNode.Value, true
				}
			}
		}
	}
	return "", false
}

// isZeroValue checks if a value represents zero (either "0" string or integer 0)
func isZeroValue(value string) bool {
	return value == "0" || value == "0s"
}
