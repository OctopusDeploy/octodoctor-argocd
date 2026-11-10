package rules

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeConfig saves content as a configuration file in a temporary directory and returns its path
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), defaultConfigFile)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestLoadIgnores(t *testing.T) {
	testCases := []struct {
		name        string
		config      string
		wantEntries int
		wantErr     string
	}{
		{"Whole rule and resources", "ignore:\n  - rule: R99\n  - rule: R19\n    resources: [\"*mesh-*\"]\n    reason: pinned\n", 2, ""},
		{"Empty file", "", 0, ""},
		{"Unknown top-level key", "ignores:\n  - rule: R99\n", 0, "field ignores not found"},
		{"Unknown entry key", "ignore:\n  - rule: R19\n    resource: [\"*\"]\n", 0, "field resource not found"},
		{"Missing rule", "ignore:\n  - reason: no rule\n", 0, `unknown or missing rule ""`},
		{"Unknown rule", "ignore:\n  - rule: R91\n", 0, `unknown or missing rule "R91"`},
		{"Empty resources", "ignore:\n  - rule: R19\n    resources: []\n", 0, "resources is empty"},
		{"Null resources", "ignore:\n  - rule: R19\n    resources:\n", 0, "resources is empty"},
		{"Whole rule twice", "ignore:\n  - rule: R99\n  - rule: R99\n", 0, "R99 is ignored as a whole more than once, or both as a whole and by resources"},
		{"Whole rule and resources for the same rule", "ignore:\n  - rule: R19\n    resources: [\"*\"]\n  - rule: R19\n", 0, "R19 is ignored as a whole more than once"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := LoadIgnores(writeConfig(t, tc.config))
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Len(t, entries, tc.wantEntries)
		})
	}

	t.Run("Missing file given with --config", func(t *testing.T) {
		_, err := LoadIgnores(filepath.Join(t.TempDir(), "missing.yaml"))
		assert.Error(t, err)
	})

	t.Run("No --config and no file in the working directory", func(t *testing.T) {
		t.Chdir(t.TempDir())
		entries, err := LoadIgnores("")
		require.NoError(t, err)
		assert.Empty(t, entries)
	})

	t.Run("No --config and a file in the working directory", func(t *testing.T) {
		t.Chdir(filepath.Dir(writeConfig(t, "ignore:\n  - rule: R99\n")))
		entries, err := LoadIgnores("")
		require.NoError(t, err)
		assert.Len(t, entries, 1)
	})
}

// fakeRule fails on the given files, or passes
type fakeRule struct {
	id    string
	files []string
	pass  bool
}

func (f fakeRule) Metadata() RuleMetadata { return RuleMetadata{ID: f.id, Severity: SeverityMedium} }

func (f fakeRule) Execute(*analysis.ArgoCDAnalysis) (RuleResult, error) {
	return RuleResult{Active: true, Passed: f.pass, Execution: analysis.RuleExecution{ID: f.id, FilesTriggered: f.files}}, nil
}

func TestExecuteAllWithIgnores(t *testing.T) {
	r19 := fakeRule{id: "R19", files: []string{"applications/mesh-a.yaml", "applications/web.yaml"}}

	testCases := []struct {
		name         string
		rule         fakeRule
		config       string
		wantFailed   []string            // files the rule still fails on, nil when it no longer fails
		wantIgnored  map[string][]string // ignored files by reason
		wantFully    int
		wantWarnings []string
	}{
		{
			name:        "Whole rule that lists no files",
			rule:        fakeRule{id: "R99"},
			config:      "ignore:\n  - rule: R99\n    reason: convention\n",
			wantIgnored: map[string][]string{"convention": nil},
			wantFully:   1,
		},
		{
			name:        "Whole rule takes every file",
			rule:        r19,
			config:      "ignore:\n  - rule: R19\n    reason: pinned\n",
			wantIgnored: map[string][]string{"pinned": r19.files},
			wantFully:   1,
		},
		{
			name:        "Some files ignored",
			rule:        r19,
			config:      "ignore:\n  - rule: R19\n    resources: [\"*mesh-*\"]\n    reason: pinned\n",
			wantFailed:  []string{"applications/web.yaml"},
			wantIgnored: map[string][]string{"pinned": {"applications/mesh-a.yaml"}},
		},
		{
			name:        "Every file ignored by resources",
			rule:        r19,
			config:      "ignore:\n  - rule: R19\n    resources: [\"*mesh-*\", \"applications/web.yaml\"]\n    reason: pinned\n",
			wantIgnored: map[string][]string{"pinned": r19.files},
			wantFully:   1,
		},
		{
			name:         "Pattern must match the whole path",
			rule:         r19,
			config:       "ignore:\n  - rule: R19\n    resources: [\"mesh-*\"]\n",
			wantFailed:   r19.files,
			wantIgnored:  map[string][]string{},
			wantWarnings: []string{`Ignore pattern "mesh-*" for R19 matched no failed file`},
		},
		{
			name:        "First matching entry claims the file",
			rule:        r19,
			config:      "ignore:\n  - rule: R19\n    resources: [\"*mesh-*\"]\n    reason: first\n  - rule: R19\n    resources: [\"*a.yaml\", \"*web*\"]\n    reason: second\n",
			wantIgnored: map[string][]string{"first": {"applications/mesh-a.yaml"}, "second": {"applications/web.yaml"}},
			wantFully:   1,
		},
		{
			name:        "Backslashes match slashes",
			rule:        fakeRule{id: "R19", files: []string{`apps\mesh-a.yaml`, "apps/mesh-b.yaml"}},
			config:      "ignore:\n  - rule: R19\n    resources: [\"apps/mesh-a.yaml\", 'apps\\mesh-b.yaml']\n",
			wantIgnored: map[string][]string{"": {`apps\mesh-a.yaml`, "apps/mesh-b.yaml"}},
			wantFully:   1,
		},
		{
			name:         "Passing rule is not ignored",
			rule:         fakeRule{id: "R14", pass: true},
			config:       "ignore:\n  - rule: R14\n",
			wantIgnored:  map[string][]string{},
			wantWarnings: []string{"Ignore entry for R14 had no effect, the rule didn't fail"},
		},
		{
			name:         "Unmatched pattern next to a matching one",
			rule:         r19,
			config:       "ignore:\n  - rule: R19\n    resources: [\"*mesh-*\", \"*typo-*\"]\n",
			wantFailed:   []string{"applications/web.yaml"},
			wantIgnored:  map[string][]string{"": {"applications/mesh-a.yaml"}},
			wantWarnings: []string{`Ignore pattern "*typo-*" for R19 matched no failed file`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ignores, err := LoadIgnores(writeConfig(t, tc.config))
			require.NoError(t, err)

			var logs bytes.Buffer
			log.SetOutput(&logs)
			t.Cleanup(func() { log.SetOutput(os.Stderr) })

			registry := &Registry{Rules: []Rule{tc.rule}, Ignores: ignores}
			result := &analysis.ArgoCDAnalysis{}
			registry.ExecuteAll(result)

			var failed []string
			if len(result.RulesFailed) > 0 {
				failed = result.RulesFailed[0].FilesTriggered
			}
			assert.Equal(t, tc.wantFailed, failed)

			ignored := map[string][]string{}
			for _, r := range result.RulesIgnored {
				ignored[r.Reason] = r.FilesTriggered
			}
			assert.Equal(t, tc.wantIgnored, ignored)
			assert.Equal(t, tc.wantFully, result.RulesFullyIgnored)
			assert.Len(t, result.RulesPassed, map[bool]int{true: 1}[tc.rule.pass])

			for _, warning := range tc.wantWarnings {
				assert.Contains(t, logs.String(), warning)
			}
			if len(tc.wantWarnings) == 0 {
				assert.NotContains(t, logs.String(), "Ignore ")
			}
		})
	}
}
