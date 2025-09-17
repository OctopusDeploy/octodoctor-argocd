package rules

import (
	"strings"

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
