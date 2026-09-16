package export

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
)

// Run exports the Argo CD resources of a live cluster into a local directory and returns the
// path of that directory, so that the regular file based analysis can run against it.
//
// An empty requestedContext means "use the current context of the kubeconfig".
func Run(requestedContext string, namespace string) (string, error) {
	contextName, err := resolveContextName(requestedContext)
	if err != nil {
		return "", err
	}

	client, err := newDynamicClient(contextName)
	if err != nil {
		return "", err
	}

	exportDirectory := exportDirectoryName(contextName)
	log.Printf("Export -> Exporting Argo CD resources from namespace %s of context %s", namespace, contextName)

	if err := prepareExportDirectory(exportDirectory); err != nil {
		return "", err
	}

	for _, kind := range exportedKinds {
		exported, err := exportKind(context.Background(), client, kind, namespace, exportDirectory)
		if err != nil {
			// A single kind we are not allowed to read, or that is not installed on this
			// cluster, should not stop the export of everything else
			if isRecoverableListError(err) {
				log.Printf("Export -> Could not list %s in namespace %s, skipping: %v", kind.name, namespace, err)
				continue
			}
			return "", fmt.Errorf("could not export %s from namespace %s: %w", kind.name, namespace, err)
		}

		log.Printf("Export -> Saved %d %s", exported, kind.name)
	}

	log.Printf("Export -> Argo CD manifests saved in %s", exportDirectory)

	return exportDirectory, nil
}

// exportKind lists a single resource type and writes every Argo CD resource it finds as its
// own YAML file. It returns how many resources were written.
func exportKind(ctx context.Context, client dynamic.Interface, kind exportedKind, namespace string, exportDirectory string) (int, error) {
	list, err := client.Resource(kind.gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	exported := 0
	for _, item := range list.Items {
		resource := item.Object

		if kind.keep != nil && !kind.keep(resource) {
			continue
		}

		pruneResource(resource)
		if kind.redact {
			redactSecretValues(resource)
		}

		path := filepath.Join(exportDirectory, kind.subfolder, item.GetName()+".yaml")
		if err := writeResource(path, resource); err != nil {
			return exported, err
		}
		exported++
	}

	return exported, nil
}

// writeResource saves a single resource as a YAML file
func writeResource(path string, resource map[string]any) error {
	contents, err := yaml.Marshal(resource)
	if err != nil {
		return fmt.Errorf("could not convert %s to YAML: %w", path, err)
	}

	if err := os.WriteFile(path, contents, 0600); err != nil {
		return fmt.Errorf("could not write %s: %w", path, err)
	}

	return nil
}

// prepareExportDirectory creates an empty export directory with one subfolder per resource
// type. An existing directory is deleted first, so that manifests left over from a previous
// export are never analysed as if they were still in the cluster.
func prepareExportDirectory(exportDirectory string) error {
	if _, err := os.Stat(exportDirectory); err == nil {
		log.Printf("Export -> Directory %s already exists, deleting it and creating it again", exportDirectory)
		if err := os.RemoveAll(exportDirectory); err != nil {
			return fmt.Errorf("could not delete existing directory %s: %w", exportDirectory, err)
		}
	}

	for _, kind := range exportedKinds {
		if err := os.MkdirAll(filepath.Join(exportDirectory, kind.subfolder), 0700); err != nil {
			return fmt.Errorf("could not create directory %s: %w", exportDirectory, err)
		}
	}

	return nil
}

// isRecoverableListError reports whether a failure to list one resource type says something
// about that type alone rather than about the cluster connection. Being denied access to
// Secrets, or asking a cluster for a CRD it does not have installed, still leaves every other
// resource type worth exporting. Anything else (the cluster being unreachable, or the
// credentials being rejected) means no resource type can be read and is not recoverable.
func isRecoverableListError(err error) bool {
	return apierrors.IsForbidden(err) || apierrors.IsNotFound(err) || apierrors.IsMethodNotSupported(err)
}
