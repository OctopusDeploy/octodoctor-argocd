package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule18(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Problematic scenario - Matrix with SCM and PR generators",
			testFile:       "appset-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Happy Path scenario - Matrix with List and PR generators",
			testFile:       "appset-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Inactive scenario - Standalone PR generator",
			testFile:       "appset-inactive.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
		{
			name:           "Inactive scenario - No PR generator at all",
			testFile:       "appset-inactive2.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r18", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R18 Failed to read test file: %s", testFilePath)

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
			rule := &Rule18{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R18 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R18 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R18 passed status mismatch")
		})
	}
}
