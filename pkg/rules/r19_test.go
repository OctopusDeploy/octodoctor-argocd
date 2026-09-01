package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule19(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Single source with targetRevision HEAD - should pass",
			testFile:       "app-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Single source with targetRevision dev - should fail",
			testFile:       "app-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Single source with chart - should be inactive",
			testFile:       "app-inactive.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
		{
			name:           "Single source without targetRevision - should pass",
			testFile:       "app-pass-default.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Multi-source with all HEAD - should pass",
			testFile:       "app-pass-multi.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Multi-source with non-HEAD targetRevision - should fail",
			testFile:       "app-fail-multi.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Multi-source with chart - should be inactive",
			testFile:       "app-inactive-multi.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r19", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R19 Failed to read test file: %s", testFilePath)

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
			rule := &Rule19{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R19 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R19 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R19 passed status mismatch")

		})
	}
}
