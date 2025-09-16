// Rule 99: Too many application sets
// This is a global rule that checks if there are too many application sets.
package rules

import "github.com/codefresh-contrib/argocd-doctor/pkg/analysis"

type Rule99 struct {
	username string
	email    string
}

const (
	ruleID          = "R99"
	ruleName        = "Check number of application sets"
	passDescription = "You have a reasonable number of application sets."
	failDescription = "You have too many application sets. This can cause performance issues and make it difficult to manage."

	appSetLimit = 10
)

func (r *Rule99) Execute(analysis *analysis.ArgoCDAnalysis) error {
	if len(analysis.ApplicationSets) > appSetLimit {
		analysis.RulesFailed = append(analysis.RulesFailed, analysis.RuleExecution{
			ID:             ruleID,
			Name:           ruleName,
			Description:    failDescription,
			FilesTriggered: []string{},
		})
		return nil
	}

	analysis.RulesPassed = append(analysis.RulesPassed, analysis.RuleExecution{
		ID:             ruleID,
		Name:           ruleName,
		Description:    passDescription,
		FilesTriggered: []string{},
	})

	return nil
}
