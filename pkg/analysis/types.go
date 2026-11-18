package analysis

type RuleExecution struct {
	ID               string
	Name             string
	Description      string
	Severity         string
	DocumentationURL string
	FilesTriggered   []string
}

// IgnoredExecution is a failed rule, or the part of it, that an ignore entry accepted
type IgnoredExecution struct {
	RuleExecution
	Reason string
}

type ArgoCDFile struct {
	FullPath string
	Name     string
	Content  string
}

type ArgoCDAnalysis struct {
	BaseDirectory string   //Base directory for calculating relative paths
	YamlFiles     []string //These are all YAML files detected
	AnalyzedFiles []string //These are Argo CD applications, application sets and configmaps

	Applications    []ArgoCDFile
	ApplicationSets []ArgoCDFile
	AppProjects     []ArgoCDFile
	ConfigMaps      []ArgoCDFile
	Secrets         []ArgoCDFile

	RulesFailed   []RuleExecution
	RulesPassed   []RuleExecution
	RulesInactive []RuleExecution

	RulesIgnored      []IgnoredExecution
	RulesFullyIgnored int // rules whose failures were all ignored
}
