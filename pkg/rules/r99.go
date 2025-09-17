// Rule 99: Too many application sets
// This is a global rule that checks if there are too many application sets.
package rules

import (
	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule99 struct {
}

const (
	appSetLimit = 10
)

func (r *Rule99) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R99",
		Name:            "Check number of application sets",
		PassDescription: "You have a reasonable number of application sets.",
		FailDescription: "You have too many application sets. This can cause performance issues and make it difficult to manage.",
		Severity:        "Medium",
	}
}

func (r *Rule99) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	if len(file_analysis.ApplicationSets) > appSetLimit {
		result := analysis.RuleExecution{
			ID:             r.Metadata().ID,
			Name:           r.Metadata().Name,
			Description:    r.Metadata().FailDescription,
			Severity:       r.Metadata().Severity,
			FilesTriggered: []string{},
		}
		return RuleResult{
			Active:    true,
			Passed:    false,
			Execution: result,
		}, nil
	}

	result := analysis.RuleExecution{
		ID:             r.Metadata().ID,
		Name:           r.Metadata().Name,
		Description:    r.Metadata().PassDescription,
		Severity:       r.Metadata().Severity,
		FilesTriggered: []string{},
	}
	return RuleResult{
		Active:    true,
		Passed:    true,
		Execution: result,
	}, nil
}
