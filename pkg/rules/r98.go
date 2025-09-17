// Rule 98: zero application sets
// This is a global rule that checks if there are no application sets.
package rules

import (
	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule98 struct {
}

func (r *Rule98) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R98",
		Name:            "Check number of application sets",
		PassDescription: "You are using application sets.",
		FailDescription: "You are not using application sets. Application Sets are a powerful feature of Argo CD that allows you to deploy applications to multiple clusters and namespaces.",
		Severity:        "Medium",
	}
}

func (r *Rule98) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	if len(file_analysis.ApplicationSets) == 0 {
		result := analysis.RuleExecution{
			ID:             r.Metadata().ID,
			Name:           r.Metadata().Name,
			Description:    r.Metadata().FailDescription,
			FilesTriggered: []string{},
		}
		return RuleResult{
			Passed:    false,
			Execution: result,
		}, nil
	}

	result := analysis.RuleExecution{
		ID:             r.Metadata().ID,
		Name:           r.Metadata().Name,
		Description:    r.Metadata().PassDescription,
		FilesTriggered: []string{},
	}
	return RuleResult{
		Passed:    true,
		Execution: result,
	}, nil
}
