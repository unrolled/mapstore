package mapstore

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Import auth for local cluster configs.
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const clusterConfigPathEnv = "MAPSTORE_CLUSTER_CONFIG_PATH"

var singleton *KubeClient

// KubeClient wires up the connection to the cluster.
type KubeClient struct {
	client    kubernetes.Interface
	ctx       context.Context
	namespace string
}

// VerifyConnection is a helper function that creates a temporary ConfigMap to ensure cluster connectivity and rbac settings.
func VerifyConnection(testMapName string, client *KubeClient) error {
	if client == nil {
		var err error

		client, err = GetKubeClient()
		if err != nil {
			return err
		}
	}

	key := "test"
	val := "ok"
	testData := map[string][]byte{key: []byte(val)}

	// Set a value.
	if err := client.Set(testMapName, testData); err != nil {
		return err
	}

	// Get a value.
	if data, err := client.Get(testMapName); err != nil {
		return err
	} else if dataVal, ok := data["test"]; !ok || string(dataVal) != val {
		return fmt.Errorf("data is mismatched")
	}

	// Delete a value.
	return client.Delete(testMapName)
}

// GetKubeClient returns the kubernetes client singleton.
func GetKubeClient() (*KubeClient, error) {
	// Try to return the singleton first.
	if singleton != nil {
		return singleton, nil
	}

	// Otherwise we need to setup the client and set the singleton.
	var err error
	var config *rest.Config

	// Determine if we should connect via our a path like `~/.kube/config`.
	if ccp, ok := os.LookupEnv(clusterConfigPathEnv); ok {
		config, err = clientcmd.BuildConfigFromFlags("", ccp)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	// Create the clientset based on the config.
	var clientset *kubernetes.Clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Grab the namespace.
	var ns string
	ns, err = getNamespace()
	if err != nil {
		return nil, err
	}

	// Set the singleton.
	singleton = &KubeClient{clientset, context.Background(), ns}

	return singleton, nil
}

func (k *KubeClient) getConfigMap(name string) (*corev1.ConfigMap, error) {
	return k.client.CoreV1().ConfigMaps(k.namespace).Get(k.ctx, name, v1.GetOptions{})
}

func (k *KubeClient) getOrCreateConfigMap(name string) (*corev1.ConfigMap, error) {
	// Attempt to fetch the existing ConfigMap.
	cm, err := k.client.CoreV1().ConfigMaps(k.namespace).Get(k.ctx, name, v1.GetOptions{})

	// If no error was returned and we have valid ConfigMap, return it.
	if err == nil && cm != nil {
		return cm, nil
	}

	// But if we had an error other than StatusReasonNotFound, return it.
	if statusError, ok := err.(*errors.StatusError); ok && statusError.Status().Reason != v1.StatusReasonNotFound {
		return nil, err
	}

	// Looks like we need to create the ConfigMap.
	cm = &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: k.namespace,
		},
	}

	return k.client.CoreV1().ConfigMaps(k.namespace).Create(k.ctx, cm, v1.CreateOptions{})
}

// Get returns the ConfigMap data.
func (k *KubeClient) Get(name string) (map[string][]byte, error) {
	cm, err := k.getConfigMap(name)
	if err != nil {
		return nil, err
	}

	return cm.BinaryData, err
}

// Set saves the ConfigMap data.
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

// Delete removes the ConfigMap entirely.
func (k *KubeClient) Delete(name string) error {
	err := k.client.CoreV1().ConfigMaps(k.namespace).Delete(k.ctx, name, v1.DeleteOptions{})

	// We can safely ignore not found errors.
	if statusError, ok := err.(*errors.StatusError); ok && statusError.Status().Reason == v1.StatusReasonNotFound {
		return nil
	}

	return err
}
