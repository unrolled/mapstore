package mapstore

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

const (
	storeTestName      = "foobar"
	storeTestNamespace = "ns-foobar"
)

func setFakeKubeClient(t *testing.T) {
	singleton = &KubeClient{fake.NewSimpleClientset(), context.Background(), storeTestNamespace}
	t.Cleanup(func() { singleton = nil })
}

func TestStoreNew(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, false)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
}

func TestStoreKeys(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	// Set a value.
	assert.NoError(t, kv.Set("k1", []byte("v1")))
	assert.NoError(t, kv.Set("k2", []byte("v2")))
	assert.NoError(t, kv.Set("k3", []byte("v3")))
	assert.NoError(t, kv.Set("k4", []byte("v4")))
	assert.NoError(t, kv.Set("k5", []byte("v5")))
	assert.Len(t, kv.internalCache, 5)

	// Should return data now.
	keys, err := kv.Keys()
	assert.NoError(t, err)
	sort.Strings(keys)

	assert.Equal(t, []string{"k1", "k2", "k3", "k4", "k5"}, keys)
}

func TestStoreGetMapData(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, false)
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

func TestStoreGetMapDataWithInternalCache(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	data, err := kv.getMapData()
	assert.NoError(t, err)
	assert.Equal(t, []byte("bar"), data["foo"])
}

func TestStoreGet(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	data, err := kv.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, []byte("bar"), data)
}

func TestStoreGetError(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar")}

	// Should return data now.
	_, err = kv.Get("this_does_not_exist")
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestStoreRaw(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
	assert.NoError(t, err)
	assert.NotNil(t, kv)
	kv.internalCache = map[string][]byte{"foo": []byte("bar"), "num": []byte("123")}

	// Should return data now.
	raw, err := kv.Raw()
	assert.NoError(t, err)
	assert.Equal(t, kv.internalCache, raw)
}

func TestStoreSet(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
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

func TestStoreSetWithInternalCache(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
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

func TestStoreSetWithDuplicateData(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
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

func TestStoreForceSet(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
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

func TestStoreForceSetWithInternalCache(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
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

func TestStoreDelete(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, false)
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
	assert.Equal(t, ErrKeyNotFound, err)
}

func TestStoreDeleteCached(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, true)
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
	assert.Equal(t, ErrKeyNotFound, err)
	assert.Len(t, kv.internalCache, 0)
}

func TestStoreTruncate(t *testing.T) {
	setFakeKubeClient(t)

	kv, err := New(storeTestName, false)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	err = kv.Set("hello", []byte("world"))
	assert.NoError(t, err)

	val, err := kv.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, []byte("world"), val)

	err = kv.Truncate()
	assert.NoError(t, err)

	_, err = kv.Get("hello")
	assert.Error(t, err)
	assert.Equal(t, ErrKeyNotFound, err)
}
