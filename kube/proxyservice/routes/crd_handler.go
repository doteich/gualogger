package routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getKubeConfig checks if the application is running inside a Kubernetes cluster or not.
// If it is running inside the cluster, it returns an in-cluster config.
// Otherwise, it returns a config from the local kubeconfig file.
func getKubeConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fallback to kubeconfig file
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	return kubeConfig.ClientConfig()
}

// GetCRDs retrieves the list of Custom Resource Definitions from the Kubernetes cluster.
func GetCRDs(c echo.Context) error {
	config, err := getKubeConfig()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get kubernetes config")
	}

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create kubernetes clientset")
	}

	crds, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list crds")
	}

	return c.JSON(http.StatusOK, crds)
}
