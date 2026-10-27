package report

import (
	"embed"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

//go:embed templates/*
var templateFS embed.FS

// templateData contains all information for the HTML report template
type templateData struct {
	Now      time.Time //when was the report was generated
	Analysis analysis.ArgoCDAnalysis
	Limit    int // max affected files shown per rule
}

var funcs = template.FuncMap{
	// sevCount counts failed rules with the given severity
	"sevCount": func(rules []analysis.RuleExecution, sev string) int {
		n := 0
		for _, r := range rules {
			if r.Severity == sev {
				n++
			}
		}
		return n
	},
	"head": func(limit int, files []string) []string {
		if len(files) > limit+1 {
			return files[:limit]
		}
		return files
	},
	"more": func(limit int, files []string) int {
		if len(files) > limit+1 {
			return len(files) - limit
		}
		return 0
	},
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// copyEmbeddedFile copies a file from the embedded filesystem to the destination path
func copyEmbeddedFile(embeddedPath string, dst string) {
	// Read the embedded file
	data, err := templateFS.ReadFile(embeddedPath)
	checkErr(err)

	// Create the destination file
	dstFile, err := os.Create(dst)
	checkErr(err)
	defer func() {
		if cerr := dstFile.Close(); cerr != nil {
			log.Fatal("close dst file: ", cerr)
		}
	}()

	// Write the data to the destination file
	_, err = dstFile.Write(data)
	checkErr(err)
}

func CreateOutputDirectory(outputDirectory string) {
	log.Println("Creating output directory at ./" + outputDirectory)
	err := os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateHTMLReport creates an HTML report from the analysis results
func GenerateHTMLReport(analysis *analysis.ArgoCDAnalysis, outputDirectory string, limit int) {
	// Copy static assets from embedded files
	copyEmbeddedFile("templates/style.css", outputDirectory+"/style.css")
	copyEmbeddedFile("templates/favicon.png", outputDirectory+"/favicon.png")
	copyEmbeddedFile("templates/octodoctor.png", outputDirectory+"/octodoctor.png")
	copyEmbeddedFile("templates/argo-cd-maintainer-badge.png", outputDirectory+"/argo-cd-maintainer-badge.png")
	copyEmbeddedFile("templates/gitops-cert-badges.png", outputDirectory+"/gitops-cert-badges.png")
	copyEmbeddedFile("templates/octopus-mark.svg", outputDirectory+"/octopus-mark.svg")

	log.Println("Rendering HTML report at ./" + outputDirectory + "/index.html")

	// Parse the HTML template from embedded files
	tmpl, err := template.New("index.html.tpl").Funcs(funcs).ParseFS(templateFS, "templates/index.html.tpl")
	if err != nil {
		log.Fatal(err)
	}

	// Create the output HTML file
	f, err := os.Create(outputDirectory + "/index.html")
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Fatal("close file: ", cerr)
		}
	}()

	// Prepare template data
	tData := templateData{
		Now:      time.Now(),
		Analysis: *analysis,
		Limit:    limit,
	}

	// Execute the template
	err = tmpl.Execute(f, tData)
	if err != nil {
		log.Fatal("Could not render HTML template: ", err)
	}
}
