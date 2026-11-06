package rules

import (
	"testing"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

func TestFailsBuild(t *testing.T) {
	high := analysis.RuleExecution{ID: "R01", Severity: SeverityHigh}
	medium := analysis.RuleExecution{ID: "R02", Severity: SeverityMedium}

	tests := []struct {
		name   string
		failed []analysis.RuleExecution
		failOn string
		want   bool
	}{
		{"No failures, medium", nil, SeverityMedium, false},
		{"No failures, high", nil, SeverityHigh, false},
		{"Medium failure, medium", []analysis.RuleExecution{medium}, SeverityMedium, true},
		{"Medium failure, high", []analysis.RuleExecution{medium}, SeverityHigh, false},
		{"Medium failure, none", []analysis.RuleExecution{medium}, FailOnNone, false},
		{"High failure, medium", []analysis.RuleExecution{high}, SeverityMedium, true},
		{"High failure, high", []analysis.RuleExecution{high}, SeverityHigh, true},
		{"High failure, none", []analysis.RuleExecution{high}, FailOnNone, false},
		{"Mixed failures, high", []analysis.RuleExecution{medium, high}, SeverityHigh, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FailsBuild(tt.failed, tt.failOn); got != tt.want {
				t.Errorf("FailsBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidFailOn(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"high", true},
		{"medium", true},
		{"none", true},
		{"low", false},
		{"", false},
		{"HIGH", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := ValidFailOn(tt.value); got != tt.want {
				t.Errorf("ValidFailOn(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
