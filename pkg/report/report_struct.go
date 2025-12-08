package report

type Metadata struct {
	Version          string `json:"version"`
	Timestamp        string `json:"timestamp"`
	ScannedDirectory string `json:"scannedDirectory"`
	Tool             string `json:"tool"`
}

type Resources struct {
	Applications    int `json:"applications"`
	ApplicationSets int `json:"applicationSets"`
	AppProjects     int `json:"appProjects"`
	ConfigMaps      int `json:"configmaps"`
	Secrets         int `json:"secrets"`
}

type Rules struct {
	Total    int `json:"total"`
	Failed   int `json:"failed"`
	Passed   int `json:"passed"`
	Inactive int `json:"inactive"`
}
type Summary struct {
	TotalYAMLFiles     int       `json:"totalYamlFiles"`
	TotalAnalyzedFiles int       `json:"totalAnalyzedFiles"`
	Resources          Resources `json:"resources"`
	Rules              Rules     `json:"rules"`
}

type Rule struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Severity      string   `json:"severity"`
	AffectedFiles []string `json:"affectedFiles"`
}

type Results struct {
	Failed   []Rule `json:"failed"`
	Passed   []Rule `json:"passed"`
	Inactive []Rule `json:"inactive"`
}
type JSONReport struct {
	Metadata Metadata `json:"metadata"`
	Summary  Summary  `json:"summary"`
	Results  Results  `json:"results"`
}
