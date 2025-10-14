package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule09(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario - Matrix with duplicate git generators",
			testFile:       "appset-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Problematic scenario - Matrix with duplicate list generators",
			testFile:       "appset-fail2.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario - Matrix with different generators",
			testFile:       "appset-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Inactive scenario - No matrix or merge generators",
			testFile:       "appset-inactive.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r09", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R09 Failed to read test file: %s", testFilePath)

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
			rule := &Rule09{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R09 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R09 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R09 passed status mismatch")

		})
	}
}
