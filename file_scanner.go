package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// discoverYamlFiles recursively scans a folder and returns a list of files with yaml or yml extensions
func (analysis *ArgoCDAnalysis) discoverYamlFiles(folderPath string) {
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

// filterArgoCDFiles filters YAML files to keep only ArgoCD Application and ApplicationSet resources
func (analysis *ArgoCDAnalysis) filterArgoCDFiles() {
	log.Printf("Filtering %d YAML files for ArgoCD resources", len(analysis.YamlFiles))

	var argoCDFiles []string

	for _, filePath := range analysis.YamlFiles {
		// Read file contents
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Scanner -> Error reading file %s: %v", filePath, err)
			continue
		}

		// Parse YAML into a map
		var yamlContent map[string]interface{}
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

		// Filter for ArgoCD resources
		if apiVersion == "argoproj.io/v1alpha1" && (kind == "Application" || kind == "ApplicationSet") {
			argoCDFiles = append(argoCDFiles, filePath)
			log.Printf("Scanner -> Found ArgoCD %s: %s", kind, filePath)
		}
	}

	log.Printf("Scanner -> Filtered to %d ArgoCD resource(s)", len(argoCDFiles))
	analysis.ArgoCDFiles = argoCDFiles
}
