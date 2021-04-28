package mapstore

import (
	"context"
	"os"
	"sync"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clusterConfigPathEnv = "CLUSTER_CONFIG_PATH"
var singleton *KubeClient
var once sync.Once

// KubeClient wires up the connection to the cluster.
type KubeClient struct {
	client    kubernetes.Interface
	ctx       context.Context
	namespace string
}

// GetKubeClient returns the kubernetes client singleton.
func GetKubeClient() (*KubeClient, error) {
	var err error

	once.Do(func() {
		var config *rest.Config

		// Determine if we should connect via our a path like `~/.kube/config`.
		if ccp, ok := os.LookupEnv(clusterConfigPathEnv); ok {
			config, err = clientcmd.BuildConfigFromFlags("", ccp)
			if err != nil {
				return
			}
		} else {
			config, err = rest.InClusterConfig()
			if err != nil {
				return
			}
		}

		// Create the clientset based on the config.
		var clientset *kubernetes.Clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return
		}

		// Grab the namespace.
		var ns string
		ns, err = getNamespace()
		if err != nil {
			return
		}

		// Create and return the client.
		singleton = &KubeClient{clientset, context.Background(), ns}
	})

	if err != nil {
		return nil, err
	}

	return singleton, nil
}

func (k *KubeClient) getConfigMap(name string) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(k.namespace).Get(k.ctx, name, v1.GetOptions{})
}

func (k *KubeClient) Get(name string) (map[string][]byte, error) {
	cm, err := k.getConfigMap(name)
	if err != nil {
		return nil, err
	}

	return cm.BinaryData, err
}

func (k *KubeClient) Set(name string, binaryData map[string][]byte) error {
	// Attempt to update if it exists.
	if cm, err := k.getConfigMap(name); err == nil {
		cm.BinaryData = binaryData
		_, updateErr := k.client.CoreV1().ConfigMaps(k.namespace).Update(k.ctx, cm, v1.UpdateOptions{})
		return updateErr
	}

	// Doesn't exists, create it instead.
	cm := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: k.namespace,
		},
		BinaryData: binaryData,
	}

	_, err := k.client.CoreV1().ConfigMaps(k.namespace).Create(k.ctx, cm, v1.CreateOptions{})

	return err
}
