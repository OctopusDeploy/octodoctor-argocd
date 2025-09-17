package analysis

type RuleExecution struct {
	ID             string
	Name           string
	Description    string
	FilesTriggered []string
}

type ArgoCDFile struct {
	FullPath string
	Name     string
	Content  string
}

type ArgoCDAnalysis struct {
	YamlFiles     []string //These are all YAML files detected
	AnalyzedFiles []string //These are Argo CD applications, application sets and configmaps

	Applications    []ArgoCDFile
	ApplicationSets []ArgoCDFile
	AppProjects     []ArgoCDFile
	ConfigMaps      []ArgoCDFile

	RulesFailed   []RuleExecution
	RulesPassed   []RuleExecution
	RulesInactive []RuleExecution
}
