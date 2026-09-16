package export

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeContextName(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		expected    string
	}{
		{
			name:        "a plain context name is left alone",
			contextName: "docker-desktop",
			expected:    "docker-desktop",
		},
		{
			name:        "dots and underscores are kept",
			contextName: "my_cluster.example",
			expected:    "my_cluster.example",
		},
		{
			name:        "an EKS ARN has its colons and slashes replaced",
			contextName: "arn:aws:eks:eu-west-1:123456789:cluster/prod",
			expected:    "arn-aws-eks-eu-west-1-123456789-cluster-prod",
		},
		{
			name:        "a GKE context keeps its underscores",
			contextName: "gke_my-project_europe-west1_my-cluster",
			expected:    "gke_my-project_europe-west1_my-cluster",
		},
		{
			name:        "spaces are replaced",
			contextName: "my cluster",
			expected:    "my-cluster",
		},
		{
			name:        "path traversal characters are replaced",
			contextName: "../../etc",
			expected:    "..-..-etc",
		},
		{
			name:        "an empty context name stays empty",
			contextName: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, sanitizeContextName(tt.contextName))
		})
	}
}

func TestExportDirectoryName(t *testing.T) {
	tests := []struct {
		name        string
		contextName string
		expected    string
	}{
		{
			name:        "a plain context name",
			contextName: "docker-desktop",
			expected:    "argocd-manifests-export-docker-desktop",
		},
		{
			name:        "an EKS ARN is sanitized first",
			contextName: "arn:aws:eks:eu-west-1:123456789:cluster/prod",
			expected:    "argocd-manifests-export-arn-aws-eks-eu-west-1-123456789-cluster-prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, exportDirectoryName(tt.contextName))
		})
	}
}
