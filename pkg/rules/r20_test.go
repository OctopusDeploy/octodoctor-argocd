package rules

import (
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/argocd-doctor/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule20(t *testing.T) {
	testCases := []struct {
		name           string
		testDir        string
		expectedPassed bool
		expectedActive bool
	}{
		{
			name:           "Pass scenario 1: Single external cluster",
			testDir:        "pass1",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Pass scenario 2: Multiple external clusters with cluster generator",
			testDir:        "pass2",
			expectedPassed: true,
			expectedActive: true,
		},
		{
			name:           "Fail scenario: Multiple external clusters without cluster generator",
			testDir:        "fail",
			expectedPassed: false,
			expectedActive: true,
		},
		{
			name:           "Inactive scenario: No external clusters",
			testDir:        "inactive",
			expectedPassed: true,
			expectedActive: false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the path to the test directory
			testDirPath := filepath.Join("..", "..", "testdata", "r20", tc.testDir)

			// Load all files from the test directory using the scanner
			testAnalysis := &analysis.ArgoCDAnalysis{}
			testAnalysis.DiscoverYamlFiles(testDirPath)
			testAnalysis.FilterArgoCDFiles()

			// Execute the rule
			rule := &Rule20{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R20 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R20 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R20 passed status mismatch")
		})
	}
}
