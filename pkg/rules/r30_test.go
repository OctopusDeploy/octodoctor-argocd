package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule30(t *testing.T) {
	testCases := []struct {
		name           string
		testFile       string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Image Updater without Git writeback",
			testFile:       "app-fail.yaml",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Image Updater with Git writeback",
			testFile:       "app-pass.yaml",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "No Image Updater annotations",
			testFile:       "app-inactive.yaml",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test file
			testFilePath := filepath.Join("..", "..", "testdata", "r30", tc.testFile)

			// Read the test file content
			content, err := os.ReadFile(testFilePath)
			require.NoError(t, err, "R30 Failed to read test file: %s", testFilePath)

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
			rule := &Rule30{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R30 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R30 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R30 passed status mismatch")

		})
	}
}
