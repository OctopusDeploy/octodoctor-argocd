package rules

import (
	"testing"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule97(t *testing.T) {
	testCases := []struct {
		name           string
		appSetCount    int
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario",
			appSetCount:    1,
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Acceptable scenario",
			appSetCount:    0,
			expectedPassed: true,
			expectedActive: false,
		},
		{
			name:           "Happy Path scenario 2",
			appSetCount:    2,
			expectedPassed: true,
			expectedActive: false,
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
			rule := &Rule97{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R97 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R97 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R97 passed status mismatch")

		})
	}
}
