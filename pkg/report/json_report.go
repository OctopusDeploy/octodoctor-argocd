package report

type Metadata struct {
	Version          string `json:"version"`
	Timestamp        string `json:"timestamp"`
	ScannedDirectory string `json:"scanned_directory"`
	Tool             string `json:"tool"`
}

type JSONReport struct {
	Metadata Metadata `json:"metadata"`
}
