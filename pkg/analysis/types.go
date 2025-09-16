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

// File rules look at at specific file. They run multiple times (one for each file detected)
type FileRule interface {
	Execute(analysis *ArgoCDAnalysis, file ArgoCDFile) error
}

// Global rules look at the entire analysis. They run only once.
type GlobalRule interface {
	Execute(analysis *ArgoCDAnalysis) error
}
