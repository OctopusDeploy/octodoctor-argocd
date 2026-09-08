package report

type Metadata struct {
	Version          string `json:"version" yaml:"version"`
	Timestamp        string `json:"timestamp" yaml:"timestamp"`
	ScannedDirectory string `json:"scannedDirectory" yaml:"scannedDirectory"`
	Tool             string `json:"tool" yaml:"tool"`
}

type Resources struct {
	Applications    int `json:"applications" yaml:"applications"`
	ApplicationSets int `json:"applicationSets" yaml:"applicationSets"`
	AppProjects     int `json:"appProjects" yaml:"appProjects"`
	ConfigMaps      int `json:"configmaps" yaml:"configmaps"`
	Secrets         int `json:"secrets" yaml:"secrets"`
}

type Rules struct {
	Total    int `json:"total" yaml:"total"`
	Failed   int `json:"failed" yaml:"failed"`
	Passed   int `json:"passed" yaml:"passed"`
	Inactive int `json:"inactive" yaml:"inactive"`
}
type Summary struct {
	TotalYAMLFiles     int       `json:"totalYamlFiles" yaml:"totalYamlFiles"`
	TotalAnalyzedFiles int       `json:"totalAnalyzedFiles" yaml:"totalAnalyzedFiles"`
	Resources          Resources `json:"resources" yaml:"resources"`
	Rules              Rules     `json:"rules" yaml:"rules"`
}

type Rule struct {
	Id               string   `json:"id" yaml:"id"`
	Name             string   `json:"name" yaml:"name"`
	Description      string   `json:"description" yaml:"description"`
	Severity         string   `json:"severity" yaml:"severity"`
	DocumentationURL string   `json:"documentationUrl" yaml:"documentationUrl"`
	AffectedFiles    []string `json:"affectedFiles" yaml:"affectedFiles"`
}

type Results struct {
	Failed   []Rule `json:"failed" yaml:"failed"`
	Passed   []Rule `json:"passed" yaml:"passed"`
	Inactive []Rule `json:"inactive" yaml:"inactive"`
}
type JSONReport struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Summary  Summary  `json:"summary" yaml:"summary"`
	Results  Results  `json:"results" yaml:"results"`
}
