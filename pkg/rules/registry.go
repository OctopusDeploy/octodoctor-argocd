package rules

import (
	"log"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type RuleResult struct {
	Passed    bool
	Execution analysis.RuleExecution
}

type RuleMetadata struct {
	ID              string
	Name            string
	PassDescription string
	FailDescription string
	Severity        string // High, Medium
}

type Rule interface {
	Execute(analysis *analysis.ArgoCDAnalysis) (RuleResult, error)
}

type Registry struct {
	Rules []Rule
}

func NewRegistry() *Registry {
	return &Registry{
		Rules: []Rule{
			&Rule01{},
			&Rule98{},
			&Rule99{},
		},
	}
}

func (r *Registry) ExecuteAll(file_analysis *analysis.ArgoCDAnalysis) error {
	// Execute all global rules at the end
	for _, rule := range r.Rules {
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
