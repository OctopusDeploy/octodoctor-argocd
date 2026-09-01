package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule02(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario",
			testFile:       "app-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Problematic scenario 2",
			testFile:       "app-fail2.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario",
			testFile:       "app-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r02", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R02 Failed to read test file: %s", testFilePath)

			// Create ArgoCDAnalysis with the test file
			testAnalysis := &analysis.ArgoCDAnalysis{
				Applications: []analysis.ArgoCDFile{
					{
						Name:     tc.testFile,
						FullPath: testFilePath,
						Content:  string(content),
					},
				},
			}

			// Execute the rule
			rule := &Rule02{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R02 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R02 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R02 passed status mismatch")

		})
	}
}
