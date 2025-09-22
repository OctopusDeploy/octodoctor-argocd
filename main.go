package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/codefresh-contrib/argocd-doctor/pkg/analysis"
	"github.com/codefresh-contrib/argocd-doctor/pkg/report"
	"github.com/codefresh-contrib/argocd-doctor/pkg/rules"
)

const (
	outputDirectory = "argocd-doctor-output"
)

func main() {
	// Parse command line flags and get target directory
	targetDir := parseDirectoryFlag()

	//Step 1 - Read all YAML files
	file_analysis := analysis.ArgoCDAnalysis{}
	// Find All YAML files (regardless of their contents)
	file_analysis.DiscoverYamlFiles(targetDir)
	file_analysis.FilterArgoCDFiles()
	log.Printf("Found %d Applications, %d ApplicationSets, %d AppProjects, %d ConfigMaps",
		len(file_analysis.Applications), len(file_analysis.ApplicationSets), len(file_analysis.AppProjects), len(file_analysis.ConfigMaps))

	//Step 2 - Execute all rules and get the results
	registry := rules.NewRegistry()
	log.Println("Checking all Argo CD files for best practices...")
	registry.ExecuteAll(&file_analysis)

	log.Printf("Analysis complete. %d rule(s) failed, %d rule(s) passed", len(file_analysis.RulesFailed), len(file_analysis.RulesPassed))

	// Step 3 - Generate a nice HTML report
	// Create output directory path relative to the target directory
	outputPath := filepath.Join(targetDir, outputDirectory)
	report.GenerateHTMLReport(&file_analysis, outputPath)

}

// parseDirectoryFlag parses command line flags and returns the target directory to scan
func parseDirectoryFlag() string {
	// Define command line flags
	var directory = flag.String("d", ".", "Directory to scan for ArgoCD files (default: current directory)")
	flag.StringVar(directory, "directory", ".", "Directory to scan for ArgoCD files (default: current directory)")
	flag.Parse()

	// Get the target directory - use flag value or current directory as fallback
	var targetDir string
	if *directory == "." {
		// Get current working directory as absolute path
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Failed to get current directory: ", err)
		}
		targetDir = currentDir
	} else {
		targetDir = *directory
	}

	// Validate that the target directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", targetDir)
	}

	return targetDir
}
