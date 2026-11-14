// Rule 08: Check if self-heal is enabled
// This rule checks if applications with an automated sync policy also have selfHeal enabled
package rules

import (
	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule08 struct {
}

func (r *Rule08) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R08",
		Name:            "Check if self-heal is enabled",
		DocSlug:         "r08-self-heal-disabled",
		PassDescription: "Applications with an automated sync policy have self-heal enabled.",
		FailDescription: "Applications have an automated sync policy without self-heal. Manual changes made directly against the cluster aren't reverted.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule08) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string

	// Iterate over all ArgoCD Application files
	for _, app := range file_analysis.Applications {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(app.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Applications without automated sync are covered by R02
		if _, hasAutomated := getYAMLValue(&yamlDoc, "spec.syncPolicy.automated"); !hasAutomated {
			continue
		}

		selfHealNode, hasSelfHeal := getYAMLValue(&yamlDoc, "spec.syncPolicy.automated.selfHeal")
		if hasSelfHeal && selfHealNode.Value == "true" {
			passedFiles = append(passedFiles, app.Name)
		} else {
			failedFiles = append(failedFiles, app.Name)
		}
	}

	// Inactive when no Application has an automated sync policy
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, nil)
}
