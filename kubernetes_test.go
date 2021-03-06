package mapstore

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	k8sTestName      = "foobar"
	k8sTestNamespace = "ns-foobar"
)

func fakeKubernetesClient() *kubeClient {
	return &kubeClient{fake.NewSimpleClientset(), context.Background(), k8sTestNamespace}
}

func TestKubernetesSingleton(t *testing.T) {
	t.Cleanup(func() { singleton = nil })
	assert.Nil(t, singleton)

	// Setup some fake values.
	os.Setenv(namespaceEnv, "foo123")
	os.Setenv(clusterConfigPathEnv, "./testdata/config.yaml")

	// Create the client.
	client, err := getKubeClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Now test that our singleton is set.
	assert.NotNil(t, singleton)
}

func TestKubernetesGetConfigMap(t *testing.T) {
	kc := fakeKubernetesClient()

	// Create a ConfigMap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(k8sTestNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: k8sTestNamespace},
		BinaryData: map[string][]byte{"foo": []byte("bar")},
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the ConfigMap.
	cm, err := kc.getConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.Equal(t, "bar", string(cm.BinaryData["foo"]))
}

func TestKubernetesGetConfigMapError(t *testing.T) {
	kc := fakeKubernetesClient()

	// Now try fetching the ConfigMap.
	cm, err := kc.getConfigMap(k8sTestName)
	assert.Error(t, err)
	assert.Nil(t, cm)
}

func TestKubernetesGet(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	// Create a ConfigMap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(k8sTestNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: k8sTestNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the ConfigMap.
	result, err := kc.get(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result)
}

func TestKubernetesGetError(t *testing.T) {
	kc := fakeKubernetesClient()

	// Now try fetching the ConfigMap.
	result, err := kc.get(k8sTestName)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestKubernetesGetOrCreateWithNoConfigMap(t *testing.T) {
	kc := fakeKubernetesClient()

	// Now try fetching the ConfigMap.
	result, err := kc.getOrCreateConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.BinaryData)

	// Create a ConfigMap that we can fetch.
	cm, err := kc.getConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.Empty(t, cm.BinaryData)
}

func TestKubernetesGetOrCreateWithExistingConfigMap(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	// Create a ConfigMap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(k8sTestNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: k8sTestNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the ConfigMap.
	result, err := kc.getOrCreateConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result.BinaryData)
}

func TestKubernetesSetCreate(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	err := kc.set(k8sTestName, data)
	assert.NoError(t, err)

	// Now try fetching the ConfigMap.
	result, err := kc.get(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result)
}

func TestKubernetesSetUpdate(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	// Create a ConfigMap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(k8sTestNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: k8sTestNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	newData := map[string][]byte{"foo": []byte("bar"), "num": []byte("one")}
	err = kc.set(k8sTestName, newData)
	assert.NoError(t, err)

	// Now try fetching the ConfigMap.
	result, err := kc.get(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newData, result)
}

func TestKubernetesDelete(t *testing.T) {
	kc := fakeKubernetesClient()

	// Create a ConfigMap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(k8sTestNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: k8sTestNamespace},
		BinaryData: map[string][]byte{"foo": []byte("bar")},
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the ConfigMap.
	err = kc.delete(k8sTestName)
	assert.NoError(t, err)
}

func TestKubernetesDeleteError(t *testing.T) {
	kc := fakeKubernetesClient()

	err := kc.delete(k8sTestName)
	assert.NoError(t, err)
}
