package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
	"github.com/OctopusDeploy/octodoctor-argocd/pkg/export"
	"github.com/OctopusDeploy/octodoctor-argocd/pkg/report"
	"github.com/OctopusDeploy/octodoctor-argocd/pkg/rules"
	"github.com/spf13/pflag"
)

var (
	// version is set at build time via ldflags
	version = "dev"
)

const (
	outputDirectory = "octodoctor-output"
)

func main() {
	// Parse command line flags
	var (
		showVersion       bool
		directory         string
		limit             int
		exportFromCluster bool
		kubeContext       string
		namespace         string
	)
	pflag.BoolVarP(&showVersion, "version", "v", false, "Print version information")
	pflag.StringVarP(&directory, "directory", "d", ".", "Directory to scan for ArgoCD files")
	pflag.IntVarP(&limit, "limit", "l", 5, "Limit number of results to display")
	pflag.BoolVarP(&exportFromCluster, "export", "e", false, "Export Argo CD manifests from a live cluster and analyze them")
	pflag.StringVarP(&kubeContext, "context", "c", "", "Kubeconfig context to export from (default: the current context)")
	pflag.StringVarP(&namespace, "namespace", "n", "argocd", "Namespace to export Argo CD manifests from")
	pflag.Parse()

	// Handle version flag
	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// Get target directory, either by exporting a live cluster or by using a local directory
	targetDir := resolveTargetDirectory(exportFromCluster, directory, kubeContext, namespace)

	//Step 1 - Read all YAML files
	file_analysis := analysis.ArgoCDAnalysis{}
	// Find All YAML files (regardless of their contents)
	file_analysis.DiscoverYamlFiles(targetDir)
	file_analysis.FilterArgoCDFiles()
	log.Printf("Found %d Applications, %d ApplicationSets, %d AppProjects, %d ConfigMaps, %d Secrets",
		len(file_analysis.Applications), len(file_analysis.ApplicationSets), len(file_analysis.AppProjects), len(file_analysis.ConfigMaps), len(file_analysis.Secrets))

	//Step 2 - Execute all rules and get the results
	registry := rules.NewRegistry()
	log.Println("Checking all Argo CD files for best practices...")
	registry.ExecuteAll(&file_analysis)

	log.Printf("Analysis complete. %d rule(s) failed, %d rule(s) passed", len(file_analysis.RulesFailed), len(file_analysis.RulesPassed))

	// Step 3 - Generate all reports (HTML, YAML, and JSON)
	// Create output directory path relative to the target directory
	outputPath := filepath.Join(targetDir, outputDirectory)
	report.CreateOutputDirectory(outputPath)
	report.GenerateStructuredReport(&file_analysis, outputPath, "yaml", version)
	report.GenerateStructuredReport(&file_analysis, outputPath, "json", version)
	report.GenerateHTMLReport(&file_analysis, outputPath, limit)

	// Step 4 - Exit with code 1 if any rules failed, otherwise 0
	if len(file_analysis.RulesFailed) > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

// resolveTargetDirectory returns the directory the analysis runs against. With --export the
// manifests of a live cluster are saved locally first and that new directory is analyzed,
// otherwise a local directory is used.
func resolveTargetDirectory(exportFromCluster bool, directory string, kubeContext string, namespace string) string {
	if !exportFromCluster {
		if pflag.CommandLine.Changed("context") || pflag.CommandLine.Changed("namespace") {
			log.Print("Ignoring --context and --namespace, they only apply together with --export")
		}
		return getTargetDirectory(directory)
	}

	if pflag.CommandLine.Changed("directory") {
		log.Fatal("Only one of --export and --directory can be used at the same time")
	}

	exportDirectory, err := export.Run(kubeContext, namespace)
	if err != nil {
		log.Fatal("Failed to export Argo CD manifests: ", err)
	}

	return exportDirectory
}

// getTargetDirectory returns the target directory to scan
func getTargetDirectory(directory string) string {
	// Get the target directory - use flag value or current directory as fallback
	var targetDir string
	if directory == "." {
		// Get current working directory as absolute path
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Failed to get current directory: ", err)
		}
		targetDir = currentDir
	} else {
		targetDir = directory
	}

	// Validate that the target directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", targetDir)
	}

	return targetDir
}
