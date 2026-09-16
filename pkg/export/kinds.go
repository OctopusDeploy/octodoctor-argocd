package export

import (
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/OctopusDeploy/octodoctor-argocd/pkg/analysis"
)

// exportedKind describes one of the resource types octodoctor knows how to analyse.
// All of them are namespaced, and every GroupVersionResource is known up front, so the export
// needs neither API discovery nor a REST mapper.
type exportedKind struct {
	name      string // human readable name, used in log lines only
	gvr       schema.GroupVersionResource
	subfolder string
	// keep decides whether a single resource is an Argo CD resource worth exporting.
	// A nil keep means every resource of this kind is exported.
	keep func(resource map[string]any) bool
	// redact replaces the values of data and stringData before writing
	redact bool
}

// exportedKinds is the complete list of resources the export pulls from a cluster. The
// ConfigMap and Secret filters call straight into pkg/analysis so that the export and the
// file scanner can never disagree about what counts as an Argo CD resource.
var exportedKinds = []exportedKind{
	{
		name:      "Applications",
		gvr:       schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"},
		subfolder: "applications",
	},
	{
		name:      "ApplicationSets",
		gvr:       schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applicationsets"},
		subfolder: "applicationsets",
	},
	{
		name:      "AppProjects",
		gvr:       schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "appprojects"},
		subfolder: "appprojects",
	},
	{
		name:      "ConfigMaps",
		gvr:       schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"},
		subfolder: "configmaps",
		keep: func(resource map[string]any) bool {
			isArgoCD, _ := analysis.IsArgoCDConfigMap(resource)
			return isArgoCD
		},
	},
	{
		name:      "Secrets",
		gvr:       schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"},
		subfolder: "secrets",
		keep: func(resource map[string]any) bool {
			isArgoCD, _ := analysis.IsArgoCDSecret(resource)
			return isArgoCD
		},
		redact: true,
	},
}
