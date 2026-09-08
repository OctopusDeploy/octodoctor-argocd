package rules

import (
	"log"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
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

// documentationBaseURL is where the mkdocs site is published. The rule pages live under
// rules/<DocSlug>/ because mkdocs uses directory URLs by default.
const documentationBaseURL = "https://octodoctor-argocd.readthedocs.io/en/latest/"

type RuleMetadata struct {
	ID              string
	Name            string
	DocSlug         string // basename of the rule page under docs/rules, without the .md extension
	PassDescription string
	FailDescription string
	Severity        string // High, Medium
}

// DocumentationURL returns the public URL of the rule's documentation page
func (m RuleMetadata) DocumentationURL() string {
	return documentationBaseURL + "rules/" + m.DocSlug + "/"
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
			&Rule07{},
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
			&Rule20{},
			&Rule22{},
			&Rule25{},
			&Rule26{},
			&Rule30{},
			&Rule50{},
			&Rule97{},
			&Rule98{},
			&Rule99{},
		},
	}
}

func (r *Registry) ExecuteAll(file_analysis *analysis.ArgoCDAnalysis) {
	// Execute all global rules at the end
	for _, rule := range r.Rules {
		result, err := rule.Execute(file_analysis)
		if err != nil {
			log.Printf("Error executing rule: %v", err)
		}
		// Set here rather than in every rule so the four rules that build their
		// RuleExecution inline cannot forget it
		result.Execution.DocumentationURL = rule.Metadata().DocumentationURL()
		if !result.Active {
			file_analysis.RulesInactive = append(file_analysis.RulesInactive, result.Execution)
		} else if result.Passed {
			file_analysis.RulesPassed = append(file_analysis.RulesPassed, result.Execution)
		} else {
			file_analysis.RulesFailed = append(file_analysis.RulesFailed, result.Execution)
		}
	}
}
