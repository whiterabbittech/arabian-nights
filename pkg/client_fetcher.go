package pkg

import (
	"github.com/whiterabbittech/arabian-nights/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type (
	clientFetcher interface {
		GetClient() (*kubernetes.Clientset, error)
	}
	outOfClusterFetcher struct{}
	inClusterFetcher    struct{}
)

var _, _ clientFetcher = &outOfClusterFetcher{}, &inClusterFetcher{}

func newOutOfClusterFetcher() outOfClusterFetcher { return outOfClusterFetcher{} }
func newInClusterFetcher() inClusterFetcher       { return inClusterFetcher{} }

func (fetcher outOfClusterFetcher) GetClient() (*kubernetes.Clientset, error) {
	var config, err = fetcher.newConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// This function builds a K8s client using the default loading rules described
// under `kubectl config --help`
func (outOfClusterFetcher) newConfig() (*rest.Config, error) {
	var loadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
	var configOverrides = &clientcmd.ConfigOverrides{}
	var kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	return kubeConfig.ClientConfig()
}

func (fetcher inClusterFetcher) GetClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func NewClient(conf *config.CLIConfig) (*kubernetes.Clientset, error) {
	var clientFetcher clientFetcher
	if conf.InCluster {
		clientFetcher = newInClusterFetcher()
	} else {
		clientFetcher = newOutOfClusterFetcher()
	}
	return clientFetcher.GetClient()
}
