package export

// redactedValue replaces the value of every Secret key. No rule ever reads a Secret value -
// R22 only looks at the key names under data - so redacting keeps every rule working while
// keeping credentials off the disk.
const redactedValue = "REDACTED-BY-OCTODOCTOR"

// prunedMetadataFields are the metadata fields that only describe the live object in the
// cluster and carry no configuration meaning. No rule reads any of them.
var prunedMetadataFields = []string{
	"uid",
	"resourceVersion",
	"generation",
	"creationTimestamp",
	"selfLink",
	"managedFields",
	"deletionTimestamp",
	"deletionGracePeriodSeconds",
	"ownerReferences",
	"finalizers",
}

const lastAppliedConfigAnnotation = "kubectl.kubernetes.io/last-applied-configuration"

// pruneResource removes the dynamic information a live cluster adds to a resource, so that the
// exported file looks like a manifest somebody would commit to Git. It modifies the resource
// in place and returns it for convenience.
func pruneResource(resource map[string]any) map[string]any {
	delete(resource, "status")

	metadata, ok := resource["metadata"].(map[string]any)
	if !ok {
		return resource
	}

	for _, field := range prunedMetadataFields {
		delete(metadata, field)
	}

	// kubectl stores a full copy of the previous manifest in this annotation, which would
	// otherwise be exported alongside the manifest itself
	if annotations, ok := metadata["annotations"].(map[string]any); ok {
		delete(annotations, lastAppliedConfigAnnotation)
		if len(annotations) == 0 {
			delete(metadata, "annotations")
		}
	}

	return resource
}

// redactSecretValues replaces every value under data and stringData while keeping the keys
// intact. It modifies the resource in place and returns it for convenience.
func redactSecretValues(resource map[string]any) map[string]any {
	for _, section := range []string{"data", "stringData"} {
		values, ok := resource[section].(map[string]any)
		if !ok {
			continue
		}
		for key := range values {
			values[key] = redactedValue
		}
	}

	return resource
}
