package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule25(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "More than two sources - should fail",
			testFile:       "app-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Two sources - should pass",
			testFile:       "app-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Single source - should be inactive",
			testFile:       "app-inactive.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r25", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R25 Failed to read test file: %s", testFilePath)

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
			rule := &Rule25{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R25 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R25 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R25 passed status mismatch")

		})
	}
}
