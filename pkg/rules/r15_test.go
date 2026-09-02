package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule15(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario",
			testFile:       "app-fail1.yaml",
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
		{
			name:           "Inactive scenario",
			testFile:       "app-inactive.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
		{
			name:           "Multi-source problematic scenario",
			testFile:       "app-fail1-multi.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Multi-source problematic scenario 2",
			testFile:       "app-fail2-multi.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Multi-source happy path",
			testFile:       "app-pass-multi.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Multi-source inactive scenario",
			testFile:       "app-inactive-multi.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r15", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R15 Failed to read test file: %s", testFilePath)

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
			rule := &Rule15{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R15 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R15 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R15 passed status mismatch")
		})
	}
}
