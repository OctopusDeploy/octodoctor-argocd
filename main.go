package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/codefresh-contrib/argocd-doctor/pkg/rules"
)

const (
	templateDirectory = "report-template/html"
	outputDirectory   = "argocd-doctor-output"
)

// Final template that contains all information. Rendered with report-template/html/index.html.tpl
type templateData struct {
	Now      time.Time //when was the report was generated
	Analysis analysis.ArgoCDAnalysis
}

func main() {

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory: ", err)
	}

	file_analysis := analysis.ArgoCDAnalysis{}
	// Find All YAML files (regardless of their contents)
	file_analysis.DiscoverYamlFiles(currentDir)
	file_analysis.FilterArgoCDFiles()
	log.Printf("Found %d Applications, %d ApplicationSets, %d AppProjects, %d ConfigMaps",
		len(file_analysis.Applications), len(file_analysis.ApplicationSets), len(file_analysis.AppProjects), len(file_analysis.ConfigMaps))

	registry := rules.NewRegistry()
	registry.ExecuteAll(&file_analysis)

	log.Printf("Found %d rules failed, %d rules passed", len(file_analysis.RulesFailed), len(file_analysis.RulesPassed))

	log.Println("Creating output directory at ./" + outputDirectory)
	err = os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		log.Fatal(err)
	}

	copy(templateDirectory+"/style.css", outputDirectory+"/style.css")
	copy(templateDirectory+"/favicon.png", outputDirectory+"/favicon.png")

	log.Println("Rendering HTML report at ./" + outputDirectory + "/index.html")

	tmpl, err := template.ParseFiles(templateDirectory + "/index.html.tpl")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(outputDirectory + "/index.html")
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Fatal("close file: ", cerr)
		}
	}()

	tData := templateData{
		Now:      time.Now(),
		Analysis: file_analysis,
	}

	err = tmpl.Execute(f, tData)

	if err != nil {
		log.Fatal("Could not render HTML template: ", err)
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(src string, dst string) {
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
