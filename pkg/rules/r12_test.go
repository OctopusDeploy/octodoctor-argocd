package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule12(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario - Merge containing Matrix (3 levels)",
			testFile:       "appset-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Problematic scenario - Matrix containing Matrix (3 levels)",
			testFile:       "appset-fail2.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario - Single Matrix (2 levels)",
			testFile:       "appset-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario - Single Merge (2 levels)",
			testFile:       "appset-pass2.yaml",
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
			testFilePath := filepath.Join("..", "..", "testdata", "r12", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R12 Failed to read test file: %s", testFilePath)

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
			rule := &Rule12{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R12 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R12 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R12 passed status mismatch")
		})
	}
}
