package rules

import (
	"strings"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"gopkg.in/yaml.v3"
)

// getYAMLValue extracts a value from YAML using a dot-separated path
func getYAMLValue(node *yaml.Node, path string) (*yaml.Node, bool) {
	if node == nil || node.Kind != yaml.DocumentNode {
		return nil, false
	}

	// Start with the root content node
	current := node.Content[0]

	// Split path by dots and navigate through the YAML structure
	parts := strings.Split(path, ".")
	for _, part := range parts {
		if current.Kind != yaml.MappingNode {
			return nil, false
		}

		// Look for the key in the mapping
		found := false
		for i := 0; i < len(current.Content); i += 2 {
			if current.Content[i].Value == part {
				current = current.Content[i+1]
				found = true
				break
			}
		}

		if !found {
			return nil, false
		}
	}

	return current, true
}

// getYAMLValueFromNode extracts a value from a YAML node using a dot-separated path
// Unlike getYAMLValue, this works on any MappingNode, not just the document root.
// This is useful when navigating within array elements or nested structures.
func getYAMLValueFromNode(node *yaml.Node, path string) (*yaml.Node, bool) {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil, false
	}

	// Start with the current node
	current := node

	// Split path by dots and navigate through the YAML structure
	parts := strings.Split(path, ".")
	for _, part := range parts {
		if current.Kind != yaml.MappingNode {
			return nil, false
		}

		// Look for the key in the mapping
		found := false
		for i := 0; i < len(current.Content); i += 2 {
			if current.Content[i].Value == part {
				current = current.Content[i+1]
				found = true
				break
			}
		}

		if !found {
			return nil, false
		}
	}

	return current, true
}

// getBoolValue safely extracts a boolean value from a YAML node
func getBoolValue(node *yaml.Node) (bool, bool) {
	if node == nil || node.Kind != yaml.ScalarNode {
		return false, false
	}

	var value bool
	if err := node.Decode(&value); err != nil {
		return false, false
	}

	return value, true
}

// determineRuleResult analyzes the collected file lists and determines the overall rule result
// This is a common function used by multiple rules to avoid code duplication
func determineRuleResult(metadata RuleMetadata, failedFiles, passedFiles, inactiveFiles []string) (RuleResult, error) {
	if len(failedFiles) > 0 {
		// Rule fails if any file has violations
		result := analysis.RuleExecution{
			ID:             metadata.ID,
			Name:           metadata.Name,
			Description:    metadata.FailDescription,
			Severity:       metadata.Severity,
			FilesTriggered: failedFiles,
		}
		return RuleResult{
			Active:    true,
			Passed:    false,
			Execution: result,
		}, nil
	}

	// Mark as inactive when no files passed the rule's relevance check — covers both
	// the case where every file was individually inactive and the case where no files
	// of the relevant type exist at all (empty lists).
	if len(passedFiles) == 0 {
		result := analysis.RuleExecution{
			ID:             metadata.ID,
			Name:           metadata.Name,
			Description:    metadata.PassDescription,
			Severity:       metadata.Severity,
			FilesTriggered: nil,
		}
		return RuleResult{
			Active:    false,
			Passed:    true,
			Execution: result,
		}, nil
	}

	// Rule passes
	allFiles := append(passedFiles, inactiveFiles...)
	result := analysis.RuleExecution{
		ID:             metadata.ID,
		Name:           metadata.Name,
		Description:    metadata.PassDescription,
		Severity:       metadata.Severity,
		FilesTriggered: allFiles,
	}
	return RuleResult{
		Active:    true,
		Passed:    true,
		Execution: result,
	}, nil
}
