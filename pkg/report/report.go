package report

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

//go:embed templates/*
var templateFS embed.FS

// templateData contains all information for the HTML report template
type templateData struct {
	Now      time.Time //when was the report was generated
	Analysis analysis.ArgoCDAnalysis
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

// GenerateHTMLReport creates an HTML report from the analysis results
func GenerateHTMLReport(analysis *analysis.ArgoCDAnalysis, outputDirectory string, limit int) {
	log.Println("Creating output directory at ./" + outputDirectory)
	err := os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Copy static assets from embedded files
	copyEmbeddedFile("templates/style.css", outputDirectory+"/style.css")
	copyEmbeddedFile("templates/favicon.png", outputDirectory+"/favicon.png")

	log.Println("Rendering HTML report at ./" + outputDirectory + "/index.html")

	// Beautify the analysis data for HTML display
	beautifyHTMLreport(analysis, limit)

	// Parse the HTML template from embedded files
	tmpl, err := template.ParseFS(templateFS, "templates/index.html.tpl")
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
	}

	// Execute the template
	err = tmpl.Execute(f, tData)
	if err != nil {
		log.Fatal("Could not render HTML template: ", err)
	}
}

// beautifyHTMLreport truncates FilesTriggered arrays in RulesFailed to maximum LIMIT entries
// If more than LIMIT files exist, the next entry becomes "+X more files" where X is the count of extra files
func beautifyHTMLreport(analysis *analysis.ArgoCDAnalysis, limit int) {
	for i := range analysis.RulesFailed {
		rule := &analysis.RulesFailed[i]
		filesCount := len(rule.FilesTriggered)

		if filesCount > limit+1 {
			// Calculate how many extra files there are
			extraFiles := filesCount - limit // We keep 3 original files + 1 summary entry

			// Truncate to first 3 files and add summary
			rule.FilesTriggered = rule.FilesTriggered[:limit]
			rule.FilesTriggered = append(rule.FilesTriggered, fmt.Sprintf("+%d more files", extraFiles))
		}
	}
}
