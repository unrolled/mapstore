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

const k8sTestName = "foobar"

func fakeKubernetesClient() *KubeClient {
	return &KubeClient{fake.NewSimpleClientset(), context.Background(), testNamespace}
}

func TestKubernetesVerifyConnection(t *testing.T) {
	err := VerifyConnection("testmap", fakeKubernetesClient())
	assert.NoError(t, err)
}

func TestKubernetesSingleton(t *testing.T) {
	t.Cleanup(func() { singleton = nil })
	assert.Nil(t, singleton)

	// Setup some fake values.
	namespaceEnv = "MAPSTORE_NAMESPACE_TEST"
	os.Setenv(namespaceEnv, "foo123")
	clusterConfigPathEnv = "MAPSTORE_CLUSTER_CONFIG_PATH_TEST"
	os.Setenv(clusterConfigPathEnv, "./testdata/config.yaml")

	// Create the client.
	client, err := GetKubeClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Now test that our singleton is set.
	assert.NotNil(t, singleton)
}

func TestKubernetesGetConfigMap(t *testing.T) {
	kc := fakeKubernetesClient()

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: testNamespace},
		BinaryData: map[string][]byte{"foo": []byte("bar")},
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	cm, err := kc.getConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.Equal(t, "bar", string(cm.BinaryData["foo"]))
}

func TestKubernetesGetConfigMapError(t *testing.T) {
	kc := fakeKubernetesClient()

	// Now try fetching the configmap.
	cm, err := kc.getConfigMap(k8sTestName)
	assert.Error(t, err)
	assert.Nil(t, cm)
}

func TestKubernetesGet(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: testNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.Get(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result)
}

func TestKubernetesGetError(t *testing.T) {
	kc := fakeKubernetesClient()

	// Now try fetching the configmap.
	result, err := kc.Get(k8sTestName)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestKubernetesGetOrCreateWithNoConfigMap(t *testing.T) {
	kc := fakeKubernetesClient()

	// Now try fetching the configmap.
	result, err := kc.getOrCreateConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.BinaryData)

	// Create a configmap that we can fetch.
	cm, err := kc.getConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.Empty(t, cm.BinaryData)
}

func TestKubernetesGetOrCreateWithExistingConfigMap(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: testNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.getOrCreateConfigMap(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result.BinaryData)
}

func TestKubernetesSetCreate(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	err := kc.Set(k8sTestName, data)
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.Get(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result)
}

func TestKubernetesSetUpdate(t *testing.T) {
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient()

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: testNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	newData := map[string][]byte{"foo": []byte("bar"), "num": []byte("one")}
	err = kc.Set(k8sTestName, newData)
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.Get(k8sTestName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newData, result)
}

func TestKubernetesDelete(t *testing.T) {
	kc := fakeKubernetesClient()

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: k8sTestName, Namespace: testNamespace},
		BinaryData: map[string][]byte{"foo": []byte("bar")},
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	err = kc.Delete(k8sTestName)
	assert.NoError(t, err)
}

func TestKubernetesDeleteError(t *testing.T) {
	kc := fakeKubernetesClient()

	err := kc.Delete(k8sTestName)
	assert.NoError(t, err)
}
