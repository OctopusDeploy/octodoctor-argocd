package export

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// resolveContextName returns the kubeconfig context to export from. An explicitly requested
// context always wins, otherwise the current context of the kubeconfig is used.
func resolveContextName(requestedContext string) (string, error) {
	if requestedContext != "" {
		return requestedContext, nil
	}

	rawConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).RawConfig()
	if err != nil {
		return "", fmt.Errorf("could not read kubeconfig: %w", err)
	}

	if rawConfig.CurrentContext == "" {
		return "", fmt.Errorf("no current context found in kubeconfig, please choose one with --context")
	}

	return rawConfig.CurrentContext, nil
}

// newDynamicClient builds a Kubernetes client for the given kubeconfig context using the
// standard kubeconfig resolution rules
func newDynamicClient(contextName string) (dynamic.Interface, error) {
	restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{CurrentContext: contextName},
	).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not connect to context %s: %w", contextName, err)
	}

	return dynamic.NewForConfig(restConfig)
}
