package rules

import (
	"testing"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule99(t *testing.T) {
	testCases := []struct {
		name           string
		appSetCount    int
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Happy Path scenario",
			appSetCount:    5,
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario at limit (10 application sets)",
			appSetCount:    10,
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Problematic scenario too application sets",
			appSetCount:    15,
			expectedPassed: false,
			expectedActive: true,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create ArgoCDAnalysis with the specified number of application sets
			testAnalysis := &analysis.ArgoCDAnalysis{
				ApplicationSets: make([]analysis.ArgoCDFile, tc.appSetCount),
			}

			// Populate with dummy application sets
			for i := 0; i < tc.appSetCount; i++ {
				testAnalysis.ApplicationSets[i] = analysis.ArgoCDFile{
					Name:     "appset",
					FullPath: "appset.yaml",
					Content:  "dummy",
				}
			}

			// Execute the rule
			rule := &Rule99{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R99 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R99 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R99 passed status mismatch")

		})
	}
}
