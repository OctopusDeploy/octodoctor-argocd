package rules

import (
	"testing"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule98(t *testing.T) {
	testCases := []struct {
		name           string
		appSetCount    int
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario zero application sets",
			appSetCount:    0,
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario at least one application set",
			appSetCount:    1,
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario many application sets",
			appSetCount:    3,
			expectedPassed: true,
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
			rule := &Rule98{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R98 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R98 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R98 passed status mismatch")

		})
	}
}
