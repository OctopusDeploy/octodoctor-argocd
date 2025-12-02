package report

type Metadata struct {
	Version          string `json:"version"`
	Timestamp        string `json:"timestamp"`
	ScannedDirectory string `json:"scanned_directory"`
	Tool             string `json:"tool"`
}

type Resources struct {
	Applications    int `json:"applications"`
	ApplicationSets int `json:"application_sets"`
	AppProjects     int `json:"app_projects"`
	ConfigMaps      int `json:"config_maps"`
	Secrets         int `json:"secrets"`
}

type Rules struct {
	Total    int `json:"total"`
	Failed   int `json:"failed"`
	Passed   int `json:"passed"`
	Inactive int `json:"inactive"`
}
type Summary struct {
	TotalYAMLFiles     int       `json:"total_yaml_files"`
	TotalAnalyzedFiles int       `json:"total_analyzed_files"`
	Resources          Resources `json:"resources"`
	Rules              Rules     `json:"rules"`
}

type Rule struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Severity      string   `json:"severity"`
	AffectedFiles []string `json:"affected_files"`
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
