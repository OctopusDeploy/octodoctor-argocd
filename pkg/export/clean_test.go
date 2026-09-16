package export

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPruneResourceRemovesStatus(t *testing.T) {
	resource := map[string]any{
		"apiVersion": "argoproj.io/v1alpha1",
		"kind":       "Application",
		"spec":       map[string]any{"project": "default"},
		"status":     map[string]any{"sync": map[string]any{"status": "Synced"}},
	}

	pruneResource(resource)

	assert.NotContains(t, resource, "status")
	assert.Contains(t, resource, "spec")
	assert.Equal(t, "argoproj.io/v1alpha1", resource["apiVersion"])
	assert.Equal(t, "Application", resource["kind"])
}

func TestPruneResourceRemovesDynamicMetadata(t *testing.T) {
	tests := []struct {
		name  string
		field string
	}{
		{name: "uid", field: "uid"},
		{name: "resourceVersion", field: "resourceVersion"},
		{name: "generation", field: "generation"},
		{name: "creationTimestamp", field: "creationTimestamp"},
		{name: "selfLink", field: "selfLink"},
		{name: "managedFields", field: "managedFields"},
		{name: "deletionTimestamp", field: "deletionTimestamp"},
		{name: "deletionGracePeriodSeconds", field: "deletionGracePeriodSeconds"},
		{name: "ownerReferences", field: "ownerReferences"},
		{name: "finalizers", field: "finalizers"},
	}

	for _, tt := range tests {
		t.Run(tt.name+" is removed", func(t *testing.T) {
			resource := map[string]any{
				"metadata": map[string]any{
					"name":   "my-app",
					tt.field: "some-value",
				},
			}

			pruneResource(resource)

			metadata, ok := resource["metadata"].(map[string]any)
			require.True(t, ok)
			assert.NotContains(t, metadata, tt.field)
			assert.Equal(t, "my-app", metadata["name"], "the name must always survive pruning")
		})
	}
}

func TestPruneResourceKeepsFieldsRulesRead(t *testing.T) {
	resource := map[string]any{
		"metadata": map[string]any{
			"name":      "my-app",
			"namespace": "argocd",
			"labels":    map[string]any{"argocd.argoproj.io/secret-type": "cluster"},
			"annotations": map[string]any{
				"argocd.argoproj.io/sync-wave": "1",
			},
			"uid": "should-go",
		},
		"spec": map[string]any{"project": "default"},
	}

	pruneResource(resource)

	metadata := resource["metadata"].(map[string]any)
	assert.Equal(t, "my-app", metadata["name"])
	assert.Equal(t, "argocd", metadata["namespace"])
	assert.Equal(t, map[string]any{"argocd.argoproj.io/secret-type": "cluster"}, metadata["labels"])
	assert.Equal(t, map[string]any{"argocd.argoproj.io/sync-wave": "1"}, metadata["annotations"])
	assert.Equal(t, map[string]any{"project": "default"}, resource["spec"])
	assert.NotContains(t, metadata, "uid")
}

func TestPruneResourceRemovesLastAppliedConfigAnnotation(t *testing.T) {
	resource := map[string]any{
		"metadata": map[string]any{
			"name": "my-app",
			"annotations": map[string]any{
				lastAppliedConfigAnnotation:    `{"apiVersion":"v1"}`,
				"argocd.argoproj.io/sync-wave": "1",
			},
		},
	}

	pruneResource(resource)

	annotations := resource["metadata"].(map[string]any)["annotations"].(map[string]any)
	assert.NotContains(t, annotations, lastAppliedConfigAnnotation)
	assert.Equal(t, "1", annotations["argocd.argoproj.io/sync-wave"])
}

func TestPruneResourceDropsAnnotationsThatBecomeEmpty(t *testing.T) {
	resource := map[string]any{
		"metadata": map[string]any{
			"name": "my-app",
			"annotations": map[string]any{
				lastAppliedConfigAnnotation: `{"apiVersion":"v1"}`,
			},
		},
	}

	pruneResource(resource)

	assert.NotContains(t, resource["metadata"].(map[string]any), "annotations")
}

func TestPruneResourceHandlesMissingSections(t *testing.T) {
	tests := []struct {
		name     string
		resource map[string]any
	}{
		{
			name:     "no metadata and no status",
			resource: map[string]any{"kind": "ConfigMap"},
		},
		{
			name:     "metadata of an unexpected type",
			resource: map[string]any{"metadata": "not-a-map"},
		},
		{
			name:     "annotations of an unexpected type",
			resource: map[string]any{"metadata": map[string]any{"annotations": "not-a-map"}},
		},
		{
			name:     "an empty resource",
			resource: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				pruneResource(tt.resource)
			})
		})
	}
}

func TestRedactSecretValues(t *testing.T) {
	resource := map[string]any{
		"kind": "Secret",
		"data": map[string]any{
			"admin.password":        "JDJhJDEwJHNvbWVoYXNo",
			"server.secretkey":      "a3ViZXJuZXRlcw==",
			"webhook.github.secret": "c3VwZXJzZWNyZXQ=",
		},
		"stringData": map[string]any{
			"token": "plaintext-token",
		},
	}

	redactSecretValues(resource)

	data := resource["data"].(map[string]any)
	// The keys must survive, because R22 looks for a data key starting with "webhook"
	assert.Len(t, data, 3)
	assert.Contains(t, data, "webhook.github.secret")
	for key, value := range data {
		assert.Equal(t, redactedValue, value, "value of %s must be redacted", key)
	}

	stringData := resource["stringData"].(map[string]any)
	assert.Equal(t, redactedValue, stringData["token"])
}

func TestRedactSecretValuesHandlesMissingSections(t *testing.T) {
	tests := []struct {
		name     string
		resource map[string]any
	}{
		{
			name:     "a secret with no data at all",
			resource: map[string]any{"kind": "Secret"},
		},
		{
			name:     "data of an unexpected type",
			resource: map[string]any{"data": "not-a-map"},
		},
		{
			name:     "empty data",
			resource: map[string]any{"data": map[string]any{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				redactSecretValues(tt.resource)
			})
		})
	}
}
