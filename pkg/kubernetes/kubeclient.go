package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// CreateClient initializes and returns a Kubernetes clientset
func CreateClient() (*kubernetes.Clientset, error) {
	var configPath string

	// Use kubeconfig from the default location (if available)
	if home := homedir.HomeDir(); home != "" {
		configPath = home + "/.kube/config"
	}

	// Build Kubernetes client from the config file
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}

	// Create a clientset using the config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
