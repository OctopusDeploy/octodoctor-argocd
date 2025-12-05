package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

func GenerateStructuredReport(analysis *analysis.ArgoCDAnalysis, outputDirectory string, language string, version string) {

	log.Println("Rendering " + language + " report at ./" + outputDirectory + "/report.json")

	// Create the output JSON file
	f, err := os.Create(outputDirectory + "/report." + language)
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Fatal("close file: ", cerr)
		}
	}()

	summary := Summary{
		TotalYAMLFiles:     len(analysis.YamlFiles),
		TotalAnalyzedFiles: len(analysis.AnalyzedFiles),
		Resources: Resources{
			Applications:    len(analysis.Applications),
			ApplicationSets: len(analysis.ApplicationSets),
			AppProjects:     len(analysis.AppProjects),
			ConfigMaps:      len(analysis.ConfigMaps),
			Secrets:         len(analysis.Secrets),
		},
		Rules: Rules{
			Total:    len(analysis.RulesFailed) + len(analysis.RulesPassed) + len(analysis.RulesInactive),
			Failed:   len(analysis.RulesFailed),
			Passed:   len(analysis.RulesPassed),
			Inactive: len(analysis.RulesInactive),
		},
	}

	// build results
	var results Results
	results.Failed = make([]Rule, len(analysis.RulesFailed))
	for i, r := range analysis.RulesFailed {
		rule := Rule{
			Id:            r.ID,
			Name:          r.Name,
			Description:   r.Description,
			Severity:      r.Severity,
			AffectedFiles: r.FilesTriggered,
		}
		results.Failed[i] = rule
	}

	results.Passed = make([]Rule, len(analysis.RulesPassed))
	for i, r := range analysis.RulesPassed {
		rule := Rule{
			Id:            r.ID,
			Name:          r.Name,
			Description:   r.Description,
			Severity:      r.Severity,
			AffectedFiles: r.FilesTriggered,
		}
		results.Passed[i] = rule
	}

	results.Inactive = make([]Rule, len(analysis.RulesInactive))
	for i, r := range analysis.RulesInactive {
		rule := Rule{
			Id:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Severity:    r.Severity,
		}
		results.Inactive[i] = rule
	}

	// build final report
	report := JSONReport{
		Metadata: Metadata{
			Version:          version,
			Timestamp:        time.Now().Format(time.RFC3339),
			ScannedDirectory: analysis.BaseDirectory,
			Tool:             "ArgoCD Doctor",
		},
		Summary: summary,
		Results: results,
	}

	// Turn structure into JSON
	if language == "json" {
		js, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			log.Fatal("marshal json: ", err)
		}
		// fmt.Printf("%s", js)
		f.WriteString(string(js))
	} else if language == "yaml" {
		yaml, err := yaml.Marshal(report)
		if err != nil {
			log.Fatal("marshal yaml: ", err)
		}
		fmt.Printf("%s", yaml)
		f.WriteString(string(yaml))
	}
	f.Sync()
}
