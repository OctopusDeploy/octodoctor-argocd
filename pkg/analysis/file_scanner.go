package analysis

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// discoverYamlFiles recursively scans a folder and returns a list of files with yaml or yml extensions
func (analysis *ArgoCDAnalysis) DiscoverYamlFiles(folderPath string) {
	// Store base directory for relative path calculations
	analysis.BaseDirectory = folderPath
	log.Printf("Scanning for YAML files under %s", folderPath)

	var yamlFiles []string
	folderCounts := make(map[string]int)

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file has a yaml or yml extension
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yaml" || ext == ".yml" {
			yamlFiles = append(yamlFiles, path)

			// Get the directory of this file
			dir := filepath.Dir(path)
			folderCounts[dir]++
		}

		return nil
	})

	if err != nil {
		// Log error but don't fail completely
		log.Printf("Scanner -> Error scanning for YAML files: %v", err)
		analysis.YamlFiles = yamlFiles
		return
	}

	// Print folder information for folders that contain YAML files in alphabetical order
	// First, collect all folder paths and sort them
	var folders []string
	for folder := range folderCounts {
		folders = append(folders, folder)
	}
	sort.Strings(folders)

	// Now print them in alphabetical order
	for _, folder := range folders {
		count := folderCounts[folder]

		// Convert absolute path to relative path
		relativePath, err := filepath.Rel(folderPath, folder)
		if err != nil {
			// If we can't get relative path, use the original folder path
			relativePath = folder
		}

		// If the relative path is ".", it means it's the root folder
		if relativePath == "." {
			relativePath = "."
		}

		log.Printf("Scanner -> [folder: %s, contains %d yaml file(s)]", relativePath, count)
	}

	log.Printf("Scanner -> Found %d YAML file(s) in %s", len(yamlFiles), folderPath)
	analysis.YamlFiles = yamlFiles
}

// filterArgoCDFiles filters YAML files to keep only ArgoCD Application, ApplicationSet resources and ConfigMaps
func (analysis *ArgoCDAnalysis) FilterArgoCDFiles() {
	log.Printf("Filtering %d YAML files for ArgoCD resources and ConfigMaps", len(analysis.YamlFiles))

	var argoCDFiles []string

	for _, filePath := range analysis.YamlFiles {
		// Read file contents
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Scanner -> Error reading file %s: %v", filePath, err)
			continue
		}

		// Parse YAML into a map
		var yamlContent map[string]any
		err = yaml.Unmarshal(fileData, &yamlContent)
		if err != nil {
			log.Printf("Scanner -> Error parsing YAML file %s: %v", filePath, err)
			continue
		}

		// Extract apiVersion and kind
		apiVersion, ok := yamlContent["apiVersion"].(string)
		if !ok {
			// apiVersion not found or not a string, skip this file
			continue
		}

		kind, ok := yamlContent["kind"].(string)
		if !ok {
			// kind not found or not a string, skip this file
			continue
		}

		// Create ArgoCDFile entry with file details
		// Calculate relative path from base directory
		relativePath, err := filepath.Rel(analysis.BaseDirectory, filePath)
		if err != nil {
			// If we can't get relative path, use the full path as fallback
			relativePath = filePath
		}
		argoCDFile := ArgoCDFile{
			FullPath: filePath,
			Name:     relativePath,
			Content:  string(fileData),
		}

		// Categorize and store the file based on its type
		if analysis.categorizeAndStoreFile(apiVersion, kind, argoCDFile, yamlContent) {
			argoCDFiles = append(argoCDFiles, filePath)
		}
	}

	log.Printf("Scanner -> Filtered to %d total ArgoCD resource(s) and ConfigMaps", len(argoCDFiles))
	analysis.AnalyzedFiles = argoCDFiles
}

// categorizeAndStoreFile categorizes and stores files according to their type
// Returns true if the file was categorized and stored, false otherwise
func (analysis *ArgoCDAnalysis) categorizeAndStoreFile(apiVersion, kind string, argoCDFile ArgoCDFile, yamlContent map[string]any) bool {
	if apiVersion == "argoproj.io/v1alpha1" && kind == "Application" {
		analysis.Applications = append(analysis.Applications, argoCDFile)
		log.Printf("Scanner -> Found ArgoCD Application: %s", argoCDFile.FullPath)
		return true
	} else if apiVersion == "argoproj.io/v1alpha1" && kind == "ApplicationSet" {
		analysis.ApplicationSets = append(analysis.ApplicationSets, argoCDFile)
		log.Printf("Scanner -> Found ArgoCD ApplicationSet: %s", argoCDFile.FullPath)
		return true
	} else if apiVersion == "argoproj.io/v1alpha1" && kind == "AppProject" {
		analysis.AppProjects = append(analysis.AppProjects, argoCDFile)
		log.Printf("Scanner -> Found ArgoCD AppProject: %s", argoCDFile.FullPath)
		return true
	} else if apiVersion == "v1" && kind == "ConfigMap" {
		isArgoCD, name := IsArgoCDConfigMap(yamlContent)

		if isArgoCD {
			analysis.ConfigMaps = append(analysis.ConfigMaps, argoCDFile)
			log.Printf("Scanner -> Found allowed ConfigMap: %s (name: %s)", argoCDFile.FullPath, name)
			return true
		}
	} else if apiVersion == "v1" && kind == "Secret" {
		isArgoCD, name := IsArgoCDSecret(yamlContent)

		if isArgoCD {
			analysis.Secrets = append(analysis.Secrets, argoCDFile)
			log.Printf("Scanner -> Found allowed Secret: %s (name: %s)", argoCDFile.FullPath, name)
			return true
		}
	}
	return false
}

// IsArgoCDConfigMap checks if a ConfigMap belong to Argo CD configuration
func IsArgoCDConfigMap(yamlContent map[string]any) (bool, string) {
	// Allowed ArgoCD ConfigMap names
	var allowedArgoCDConfigMaps = []string{
		"argocd-cm",
		"argocd-cmd-params-cm",
		"argocd-rbac-cm",
		"argocd-tls-certs-cm",
		"argocd-ssh-known-hosts-cm",
	}

	// Read metadata.name
	metadata, ok := yamlContent["metadata"].(map[string]any)
	if !ok {
		// metadata not found or not a map
		return false, ""
	}

	name, ok := metadata["name"].(string)
	if !ok {
		// name not found or not a string
		return false, ""
	}

	// Check if the ConfigMap name is in the allowed list
	if slices.Contains(allowedArgoCDConfigMaps, name) {
		return true, name
	}

	return false, name
}

// IsArgoCDSecret checks if a Secret belongs to Argo CD configuration
func IsArgoCDSecret(yamlContent map[string]any) (bool, string) {
	// Allowed ArgoCD Secret names
	var allowedArgoCDSecrets = []string{
		"argocd-secret",
	}

	// Read metadata.name
	metadata, ok := yamlContent["metadata"].(map[string]any)
	if !ok {
		// metadata not found or not a map
		return false, ""
	}

	name, ok := metadata["name"].(string)
	if !ok {
		// name not found or not a string
		return false, ""
	}

	// Check if the Secret name is in the allowed list
	if slices.Contains(allowedArgoCDSecrets, name) {
		return true, name
	}

	// Check if this is a cluster type secret (for rules that deal with external clusters)
	// Look for the label argocd.argoproj.io/secret-type: cluster
	labels, ok := metadata["labels"].(map[string]any)
	if ok {
		secretType, ok := labels["argocd.argoproj.io/secret-type"].(string)
		if ok && secretType == "cluster" {
			return true, name
		}
	}

	return false, name
}
