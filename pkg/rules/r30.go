// Rule 30: Check if Image Updater uses Git writeback
// This rule checks if applications using Argo CD Image Updater have Git writeback configured
package rules

import (
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

type Rule30 struct {
}

const (
	imageUpdaterPrefix      = "argocd-image-updater.argoproj.io"
	writeBackMethodKey      = "argocd-image-updater.argoproj.io/write-back-method"
	writeBackMethodGitValue = "git"
)

func (r *Rule30) Metadata() RuleMetadata {
	return RuleMetadata{
		ID:              "R30",
		Name:            "Check if Image Updater uses Git writeback",
		DocSlug:         "r30-image-updater-without-git-writeback",
		PassDescription: "Applications using Image Updater have Git writeback configured.",
		FailDescription: "Applications using Image Updater without Git writeback. This is against GitOps principles.",
		Severity:        SeverityHigh,
	}
}

func (r *Rule30) Execute(file_analysis *analysis.ArgoCDAnalysis) (RuleResult, error) {
	var failedFiles []string
	var passedFiles []string
	var inactiveFiles []string

	// Iterate over all ArgoCD Application files
	for _, app := range file_analysis.Applications {
		// Parse YAML content using Node API for path-based access
		var yamlDoc yaml.Node
		err := yaml.Unmarshal([]byte(app.Content), &yamlDoc)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Check if annotations exist
		annotationsNode, hasAnnotations := getYAMLValue(&yamlDoc, "metadata.annotations")
		if !hasAnnotations || annotationsNode.Kind != yaml.MappingNode {
			// No annotations, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
			continue
		}

		// Check for image-updater annotations
		hasImageUpdater := false
		hasGitWriteBack := false

		// Iterate through all annotations
		for i := 0; i < len(annotationsNode.Content); i += 2 {
			annotationKey := annotationsNode.Content[i].Value
			annotationValue := annotationsNode.Content[i+1].Value

			// Check if this is an image-updater annotation
			if strings.HasPrefix(annotationKey, imageUpdaterPrefix) {
				hasImageUpdater = true

				// Check if this is the write-back-method annotation with value "git"
				if annotationKey == writeBackMethodKey && annotationValue == writeBackMethodGitValue {
					hasGitWriteBack = true
				}
			}
		}

		// Determine result based on image updater usage
		if !hasImageUpdater {
			// No image-updater annotations, mark as inactive
			inactiveFiles = append(inactiveFiles, app.Name)
		} else if hasImageUpdater && !hasGitWriteBack {
			// Has image-updater but no Git writeback, rule fails
			failedFiles = append(failedFiles, app.Name)
		} else {
			// Has image-updater with Git writeback, rule passes
			passedFiles = append(passedFiles, app.Name)
		}
	}

	// Determine overall rule result
	return determineRuleResult(r.Metadata(), failedFiles, passedFiles, inactiveFiles)
}
