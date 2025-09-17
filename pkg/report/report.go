package report

import (
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
)

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

func Copy(src string, dst string) {
	// Open the source file
	srcFile, err := os.Open(src)
	checkErr(err)
	defer func() {
		if cerr := srcFile.Close(); cerr != nil {
			log.Fatal("close src file: ", cerr)
		}
	}()

	// Create the destination file
	dstFile, err := os.Create(dst)
	checkErr(err)
	defer func() {
		if cerr := dstFile.Close(); cerr != nil {
			log.Fatal("close dst file: ", cerr)
		}
	}()

	// Copy the contents from src to dst
	_, err = io.Copy(dstFile, srcFile)
	checkErr(err)
}

// GenerateHTMLReport creates an HTML report from the analysis results
func GenerateHTMLReport(analysis *analysis.ArgoCDAnalysis, templateDirectory, outputDirectory string) {
	log.Println("Creating output directory at ./" + outputDirectory)
	err := os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Copy static assets
	Copy(templateDirectory+"/style.css", outputDirectory+"/style.css")
	Copy(templateDirectory+"/favicon.png", outputDirectory+"/favicon.png")

	log.Println("Rendering HTML report at ./" + outputDirectory + "/index.html")

	// Parse the HTML template
	tmpl, err := template.ParseFiles(templateDirectory + "/index.html.tpl")
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
