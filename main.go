package main

import (
	"log"
	"os"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/codefresh-contrib/argocd-doctor/pkg/report"
	"github.com/codefresh-contrib/argocd-doctor/pkg/rules"
)

const (
	templateDirectory = "report-template/html"
	outputDirectory   = "argocd-doctor-output"
)

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

	// Generate HTML report
	report.GenerateHTMLReport(&file_analysis, templateDirectory, outputDirectory)

}
