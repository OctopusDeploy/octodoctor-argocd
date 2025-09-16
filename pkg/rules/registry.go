package rules

import (
	"log"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type Registry struct {
	globalRules []analysis.GlobalRule
	fileRules   []analysis.FileRule
}

func NewRegistry() *Registry {
	return &Registry{
		globalRules: []analysis.GlobalRule{
			&Rule98{},
			&Rule99{},
		},
		fileRules: []analysis.FileRule{
			// Add file rules here
		},
	}
}

func (r *Registry) ExecuteAll(file_analysis *analysis.ArgoCDAnalysis) error {
	// Execute all global rules at the end
	for _, rule := range r.globalRules {
		result, err := rule.Execute(file_analysis)
		if err != nil {
			log.Printf("Error executing rule: %v", err)
		}
		if result.Passed {
			file_analysis.RulesPassed = append(file_analysis.RulesPassed, result.Execution)
		} else {
			file_analysis.RulesFailed = append(file_analysis.RulesFailed, result.Execution)
		}
	}
	return nil
}
