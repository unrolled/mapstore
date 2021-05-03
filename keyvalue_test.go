package mapstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

func setKeyValueFakeKubeClient(t *testing.T, ns string) {
	singleton = &KubeClient{fake.NewSimpleClientset(), context.Background(), ns}
	t.Cleanup(func() { singleton = nil })
}

func TestKeyValueNew(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, false)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
}

func TestKeyValueGetMapData(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, false)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	// Should return nothing.
	data, err := kv.getMapData()
	assert.NoError(t, err)
	assert.Empty(t, data)

	// Set a value.
	err = kv.Set("foo", []byte("bar"))
	assert.NoError(t, err)

	// Should return data now.
	data, err = kv.getMapData()
	assert.NoError(t, err)
	assert.Equal(t, []byte("bar"), data["foo"])
}

func TestKeyValueGetMapDataWithInternalCache(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	data, err := kv.getMapData()
	assert.NoError(t, err)
	assert.Equal(t, []byte("bar"), data["foo"])
}

func TestKeyValueGet(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	data, err := kv.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, []byte("bar"), data)
}

func TestKeyValueGetError(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	_, err = kv.Get("this_does_not_exist")
	assert.Equal(t, ErrKeyValueNotFound, err)
}

func TestKeyValueSet(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	// Should return data now.
	err = kv.Set("hello", []byte("world"))
	assert.NoError(t, err)

	// Now get the value.
	data, err := kv.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), data)
}

func TestKeyValueSetWithInternalCache(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	err = kv.Set("hello", []byte("world"))
	assert.NoError(t, err)
	assert.Len(t, kv.internalCache, 2)

	// Now get the value.
	data, err := kv.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), data)
}

func TestKeyValueSetWithDuplicateData(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Explicitly set the client to nil so we won't make any k8s calls.
	kv.client = nil

	// Should return data now.
	err = kv.Set("foo", []byte("bar"))
	assert.NoError(t, err)
	assert.Len(t, kv.internalCache, 1)

	// Now get the value.
	data, err := kv.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, []byte("bar"), data)
}

func TestKeyValueForceSet(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	// Should return data now.
	err = kv.ForceSet("hello", []byte("world"))
	assert.NoError(t, err)

	// Now get the value.
	data, err := kv.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), data)
}

func TestKeyValueForceSetWithInternalCache(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	err = kv.ForceSet("hello", []byte("world"))
	assert.NoError(t, err)
	assert.Len(t, kv.internalCache, 2)

	// Now get the value.
	data, err := kv.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), data)
}

func TestKeyValueDelete(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, false)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	err = kv.Set("hello", []byte("world"))
	assert.NoError(t, err)

	val, err := kv.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), val)

	err = kv.Delete("hello")
	assert.NoError(t, err)

	_, err = kv.Get("hello")
	assert.Error(t, err)
	assert.Equal(t, ErrKeyValueNotFound, err)
}

func TestKeyValueDeleteCached(t *testing.T) {
	setKeyValueFakeKubeClient(t, testNamespace)

	name := "testmap"
	kv, err := NewKeyValue(name, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	err = kv.Set("hello", []byte("world"))
	assert.NoError(t, err)
	assert.Len(t, kv.internalCache, 1)

	err = kv.Delete("hello")
	assert.NoError(t, err)
	assert.Len(t, kv.internalCache, 0)

	_, err = kv.Get("hello")
	assert.Error(t, err)
	assert.Equal(t, ErrKeyValueNotFound, err)
	assert.Len(t, kv.internalCache, 0)
}
