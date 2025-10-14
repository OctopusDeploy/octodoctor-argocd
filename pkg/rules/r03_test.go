package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule03(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario - requeueAfterSeconds too low (20 seconds)",
			testFile:       "appset-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Problematic scenario - requeueAfterSeconds too low (120 seconds)",
			testFile:       "appset-fail2.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario - requeueAfterSeconds high enough (1800 seconds)",
			testFile:       "appset-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario - requeueAfterSeconds high enough (3000 seconds in nested generator)",
			testFile:       "appset-pass2.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario - requeueAfterSeconds high enough (600 seconds in nested generator)",
			testFile:       "appset-pass3.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Inactive scenario - No requeueAfterSeconds defined",
			testFile:       "appset-inactive.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r03", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R03 Failed to read test file: %s", testFilePath)

			// Create ArgoCDAnalysis with the test file
			testAnalysis := &analysis.ArgoCDAnalysis{
				ApplicationSets: []analysis.ArgoCDFile{
					{
						Name:     tc.testFile,
						FullPath: testFilePath,
						Content:  string(content),
					},
				},
			}

			// Execute the rule
			rule := &Rule03{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R03 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R03 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R03 passed status mismatch")

		})
	}
}
