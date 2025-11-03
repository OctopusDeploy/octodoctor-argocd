package rules

import (
	"path/filepath"
	"testing"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule07(t *testing.T) {
	testCases := []struct {
		name           string
		testDir        string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Inactive scenario: timeout.reconciliation not present",
			testDir:        "inactive",
			expectedPassed: true,
			expectedActive: false,
		},
		{
			name:           "Pass scenario 1: timeout.reconciliation is not 0",
			testDir:        "pass1",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Pass scenario 2: timeout.reconciliation is 0 but cache expiration is set",
			testDir:        "pass2",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Fail scenario: timeout.reconciliation is 0 and cache expiration is default",
			testDir:        "fail",
			expectedPassed: false,
			expectedActive: true,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test directory
			testDirPath := filepath.Join("..", "..", "testdata", "r07", tc.testDir)

			// Load all files from the test directory using the scanner
			testAnalysis := &analysis.ArgoCDAnalysis{}
			testAnalysis.DiscoverYamlFiles(testDirPath)
			testAnalysis.FilterArgoCDFiles()

			// Execute the rule
			rule := &Rule07{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R07 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R07 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R07 passed status mismatch")
		})
	}
}
