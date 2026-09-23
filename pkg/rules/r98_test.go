package rules

import (
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRule98(t *testing.T) {
	testCases := []struct {
		name                string
		testDir             string
		expectedAppCount    int
		expectedAppSetCount int
		expectedPassed      bool
		expectedActive      bool
	}{
		{
			name:                "Applications without application sets - should fail",
			testDir:             "apps-no-appsets",
			expectedAppCount:    2,
			expectedAppSetCount: 0,
			expectedPassed:      false,
			expectedActive:      true,
		},
		{
			name:                "No Applications and no application sets - should be inactive",
			testDir:             "empty",
			expectedAppCount:    0,
			expectedAppSetCount: 0,
			expectedPassed:      true,
			expectedActive:      false,
		},
		{
			name:                "Application sets without Applications - should pass",
			testDir:             "appsets-only",
			expectedAppCount:    0,
			expectedAppSetCount: 1,
			expectedPassed:      true,
			expectedActive:      true,
		},
		{
			name:                "Applications and application sets - should pass",
			testDir:             "apps-and-appsets",
			expectedAppCount:    1,
			expectedAppSetCount: 1,
			expectedPassed:      true,
			expectedActive:      true,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Scan the test directory the same way the CLI does, since this is a global rule
			testAnalysis := &analysis.ArgoCDAnalysis{}
			testAnalysis.DiscoverYamlFiles(filepath.Join("..", "..", "testdata", "r98", tc.testDir))
			testAnalysis.FilterArgoCDFiles()

			require.Len(t, testAnalysis.Applications, tc.expectedAppCount, "R98 unexpected number of Applications in %s", tc.testDir)
			require.Len(t, testAnalysis.ApplicationSets, tc.expectedAppSetCount, "R98 unexpected number of application sets in %s", tc.testDir)

			// Execute the rule
			rule := &Rule98{}
			result, err := rule.Execute(testAnalysis)

			require.NoError(t, err, "Rule R98 could not run")

			// Assert the rule active status
			assert.Equal(t, tc.expectedActive, result.Active, "Rule R98 active status mismatch")
			assert.Equal(t, tc.expectedPassed, result.Passed, "Rule R98 passed status mismatch")
		})
	}
}
