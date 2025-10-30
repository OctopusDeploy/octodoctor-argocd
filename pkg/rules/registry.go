package rules

import (
	"log"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

type RuleResult struct {
	Passed    bool
	Active    bool
	Execution analysis.RuleExecution
}

// Severity constants for rule metadata
const (
	SeverityHigh   = "high"
	SeverityMedium = "medium"
	SeverityLow    = "low"
)

type RuleMetadata struct {
	ID              string
	Name            string
	PassDescription string
	FailDescription string
	Severity        string // High, Medium
}

type Rule interface {
	Metadata() RuleMetadata
	Execute(analysis *analysis.ArgoCDAnalysis) (RuleResult, error)
}

type Registry struct {
	Rules []Rule
}

func NewRegistry() *Registry {
	return &Registry{
		Rules: []Rule{
			&Rule01{},
			&Rule02{},
			&Rule03{},
			&Rule04{},
			&Rule05{},
			&Rule06{},
			&Rule09{},
			&Rule10{},
			&Rule11{},
			&Rule12{},
			&Rule14{},
			&Rule15{},
			&Rule16{},
			&Rule17{},
			&Rule18{},
			&Rule19{},
			&Rule22{},
			&Rule30{},
			&Rule97{},
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
		if !result.Active {
			file_analysis.RulesInactive = append(file_analysis.RulesInactive, result.Execution)
		} else if result.Passed {
			file_analysis.RulesPassed = append(file_analysis.RulesPassed, result.Execution)
		} else {
			file_analysis.RulesFailed = append(file_analysis.RulesFailed, result.Execution)
		}
	}
	return nil
}
