package clienthelper

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientset offers a convenience factory for Kubernetes clients,
// depending on whether a path to a kubernetes cluster configuration is
// provided or not.
func NewClientset(kubeconfig string) (*kubernetes.Clientset, error) {

	if kubeconfig == "" {
		// We consider that we are running in-cluster as long
		// as there is no kubeconfig provided
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		// creates the clientset
		return kubernetes.NewForConfig(config)
	}

	// we are running outside a cluster.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	return kubernetes.NewForConfig(config)

}
