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

func fakeKubernetesClient(ns string) *KubeClient {
	return &KubeClient{fake.NewSimpleClientset(), context.Background(), ns}
}

func TestKubernetesVerifyConnection(t *testing.T) {
	err := VerifyConnection("testmap", fakeKubernetesClient(testNamespace))
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
	name := "foobar"
	kc := fakeKubernetesClient(testNamespace)

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: testNamespace},
		BinaryData: map[string][]byte{"foo": []byte("bar")},
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	cm, err := kc.getConfigMap(name)
	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.Equal(t, "bar", string(cm.BinaryData["foo"]))
}

func TestKubernetesGetConfigMapError(t *testing.T) {
	name := "foobar"
	kc := fakeKubernetesClient(testNamespace)

	// Now try fetching the configmap.
	cm, err := kc.getConfigMap(name)
	assert.Error(t, err)
	assert.Nil(t, cm)
}

func TestKubernetesGet(t *testing.T) {
	name := "foobar"
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient(testNamespace)

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: testNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.Get(name)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result)
}

func TestKubernetesGetError(t *testing.T) {
	name := "foobar"
	kc := fakeKubernetesClient(testNamespace)

	// Now try fetching the configmap.
	result, err := kc.Get(name)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestKubernetesGetOrCreateWithNoConfigMap(t *testing.T) {
	name := "foobar"
	kc := fakeKubernetesClient(testNamespace)

	// Now try fetching the configmap.
	result, err := kc.getOrCreateConfigMap(name)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.BinaryData)

	// Create a configmap that we can fetch.
	cm, err := kc.getConfigMap(name)
	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.Empty(t, cm.BinaryData)
}

func TestKubernetesGetOrCreateWithExistingConfigMap(t *testing.T) {
	name := "foobar"
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient(testNamespace)

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: testNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.getOrCreateConfigMap(name)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result.BinaryData)
}

func TestKubernetesSetCreate(t *testing.T) {
	name := "foobar"
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient(testNamespace)

	err := kc.Set(name, data)
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.Get(name)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, data, result)
}

func TestKubernetesSetUpdate(t *testing.T) {
	name := "foobar"
	data := map[string][]byte{"foo": []byte("bar")}
	kc := fakeKubernetesClient(testNamespace)

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: testNamespace},
		BinaryData: data,
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	newData := map[string][]byte{"foo": []byte("bar"), "num": []byte("one")}
	err = kc.Set(name, newData)
	assert.NoError(t, err)

	// Now try fetching the configmap.
	result, err := kc.Get(name)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newData, result)
}

func TestKubernetesDelete(t *testing.T) {
	name := "foobar"
	kc := fakeKubernetesClient(testNamespace)

	// Create a configmap that we can fetch.
	_, err := kc.client.CoreV1().ConfigMaps(testNamespace).Create(kc.ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: testNamespace},
		BinaryData: map[string][]byte{"foo": []byte("bar")},
	}, v1.CreateOptions{})
	assert.NoError(t, err)

	// Now try fetching the configmap.
	err = kc.Delete(name)
	assert.NoError(t, err)
}

func TestKubernetesDeleteError(t *testing.T) {
	name := "foobar"
	kc := fakeKubernetesClient(testNamespace)

	err := kc.Delete(name)
	assert.NoError(t, err)
}
