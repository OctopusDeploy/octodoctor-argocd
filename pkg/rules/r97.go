// Rule 97: exactly one application set
// This is a global rule that checks if there is exactly one application set.
package rules

import (
	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule97 struct {
}

func (r *Rule97) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R97",
		Name:            "Check for exactly one application set",
		PassDescription: "You have a reasonable number of application sets.",
		FailDescription: "You have exactly one application set. This usually results in a very complex YAML file.Consider using multiple application sets",
		Severity:        "medium",
	}
}

func (r *Rule97) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	appSetCount := len(file_analysis.ApplicationSets)

	// Rule is inactive if 0 or more than 1 application sets
	if appSetCount == 0 || appSetCount > 1 {
		result := analysis.RuleExecution{
			ID:             r.Metadata().ID,
			Name:           r.Metadata().Name,
			Description:    r.Metadata().PassDescription,
			Severity:       r.Metadata().Severity,
			FilesTriggered: []string{},
		}
		return RuleResult{
			Active:    false,
			Passed:    true,
			Execution: result,
		}, nil
	}

	// Rule fails if exactly 1 application set
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
