// Rule 98: zero application sets
// This is a global rule that checks if there are Applications but no application sets.
package rules

import (
	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule98 struct {
}

func (r *Rule98) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R98",
		Name:            "Check number of application sets",
		DocSlug:         "r98-no-application-sets",
		PassDescription: "You are using application sets.",
		FailDescription: "You are not using application sets. Application Sets are a powerful feature of Argo CD that allows you to deploy applications to multiple clusters and namespaces.",
		Severity:        SeverityMedium,
	}
}

func (r *Rule98) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	appCount := len(file_analysis.Applications)
	appSetCount := len(file_analysis.ApplicationSets)

	// Rule is inactive if there is nothing to generate: no Applications and no application sets
	if appCount == 0 && appSetCount == 0 {
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

	// Rule fails if there are Applications but no application sets
	if appSetCount == 0 {
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
