package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

func GenerateJSONReport(analysis *analysis.ArgoCDAnalysis, outputDirectory string, version string) {
	// log.Println("Creating output directory at ./" + outputDirectory)
	// err := os.MkdirAll(outputDirectory, 0755)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Rendering JSON report at ./" + outputDirectory + "/report.json")

	// Create the output JSON file
	f, err := os.Create(outputDirectory + "/report.json")
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Fatal("close file: ", cerr)
		}
	}()

	metadata := Metadata{
		Version:          version,
		Timestamp:        time.Now().Format(time.RFC3339),
		ScannedDirectory: analysis.BaseDirectory,
		Tool:             "ArgoCD Doctor",
	}

	report := JSONReport{
		Metadata: metadata,
	}

	// Turn structure into JSON
	js, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatal("marshal json: ", err)
	}
	fmt.Printf("%s", js)
	f.WriteString(string(js))
	f.Sync()
}
