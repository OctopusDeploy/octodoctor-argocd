// Rule 99: Too many application sets
// This is a global rule that checks if there are too many application sets.
package rules

import (
	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Rule99 struct {
	ruleID   string
	ruleName string
}

const (
	ID              = "R99"
	Name            = "Check number of application sets"
	passDescription = "You have a reasonable number of application sets."
	failDescription = "You have too many application sets. This can cause performance issues and make it difficult to manage."

	appSetLimit = 10
)

func (r *Rule99) Execute(file_analysis *analysis.ArgoCDAnalysis) error {
	if len(file_analysis.ApplicationSets) > appSetLimit {
		file_analysis.RulesFailed = append(file_analysis.RulesFailed, analysis.RuleExecution{
			ID:             ID,
			Name:           Name,
			Description:    failDescription,
			FilesTriggered: []string{},
		})
		return nil
	}

	file_analysis.RulesPassed = append(file_analysis.RulesPassed, analysis.RuleExecution{
		ID:             ID,
		Name:           Name,
		Description:    passDescription,
		FilesTriggered: []string{},
	})

	return nil
}
