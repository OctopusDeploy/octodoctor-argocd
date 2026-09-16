package export

import "regexp"

const exportDirectoryPrefix = "argocd-manifests-export-"

// unsafeFileNameChars matches every character that is not safe to use in a directory name.
// Kubeconfig context names frequently contain characters that are illegal or awkward in a
// path, for example an EKS context is usually an ARN such as
// arn:aws:eks:eu-west-1:123456789:cluster/prod
var unsafeFileNameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

// sanitizeContextName replaces every character that is unsafe in a directory name with a dash
func sanitizeContextName(contextName string) string {
	return unsafeFileNameChars.ReplaceAllString(contextName, "-")
}

// exportDirectoryName returns the name of the directory that holds the manifests exported
// from the given kubeconfig context
func exportDirectoryName(contextName string) string {
	return exportDirectoryPrefix + sanitizeContextName(contextName)
}
