package mapstore

import (
	"bytes"
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ErrKeyValueNotFound is returned when looking up a value that does not exist.
var ErrKeyValueNotFound = fmt.Errorf("key was not found")

// SimpleInterface defines the required methods to satisfy the KeyValue implementation.
type SimpleInterface interface {
	Keys() ([]string, error)
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	Reset() error
}

// ComplexInterface defines the required methods and a few optional methods for the KeyValue implementation.
type ComplexInterface interface {
	SimpleInterface
	Raw() (map[string][]byte, error)
	ForceSet(key string, value []byte) error
}

// Verify we meet the requirements for our own internfaces.
var _ SimpleInterface = &KeyValue{}
var _ ComplexInterface = &KeyValue{}

// KeyValue is a thread safe key value store backed by a Kubernetes ConfigMap.
type KeyValue struct {
	*sync.RWMutex
	configmapName string
	client        *KubeClient
	cacheEnabled  bool
	internalCache map[string][]byte
}

// NewKeyValue returns a newly setup instance of KeyValue.
func NewKeyValue(cmName string, cacheInternally bool) (*KeyValue, error) {
	// Grab the KubeClient.
	kubeClient, err := GetKubeClient()
	if err != nil {
		return nil, err
	}

	// If we are caching internally, fetch the data first.
	cache := map[string][]byte{}
	if cacheInternally {
		if cm, err := kubeClient.getOrCreateConfigMap(cmName); err != nil {
			return nil, err
		} else if cm.BinaryData != nil {
			cache = cm.BinaryData
		}
	}

	return &KeyValue{
		RWMutex:       &sync.RWMutex{},
		configmapName: cmName,
		client:        kubeClient,
		cacheEnabled:  cacheInternally,
		internalCache: cache,
	}, nil
}

func (k *KeyValue) getMapData() (map[string][]byte, error) {
	if k.cacheEnabled {
		return k.internalCache, nil
	}

	data, err := k.client.Get(k.configmapName)

	// Determine if the error was a "not found" error or not.
	statusError, statusCastOk := err.(*errors.StatusError)
	isNotFound := statusCastOk && statusError.Status().Reason == v1.StatusReasonNotFound
	if err != nil && !isNotFound {
		return nil, err
	}

	if data == nil {
		data = map[string][]byte{}
	}

	return data, nil
}

func (k *KeyValue) Keys() ([]string, error) {
	k.RLock()
	defer k.RUnlock()

	// Grab the data map.
	dataMap, err := k.getMapData()
	if err != nil {
		return nil, err
	}

	// Lookup all the keys.
	keys := make([]string, 0, len(dataMap))
	for k := range dataMap {
		keys = append(keys, k)
	}

	return keys, nil
}

// Get uses the supplied key and attempts to return the coorsponding value from the ConfigMap.
func (k *KeyValue) Get(key string) ([]byte, error) {
	k.RLock()
	defer k.RUnlock()

	// Grab the data map.
	dataMap, err := k.getMapData()
	if err != nil {
		return nil, err
	}

	// Lookup the value, return not found if it failed.
	val, ok := dataMap[key]
	if !ok {
		return nil, ErrKeyValueNotFound
	}

	return val, nil
}

// Raw returns the actual underlying map data.
func (k *KeyValue) Raw() (map[string][]byte, error) {
	k.RLock()
	defer k.RUnlock()

	// Grab the data map.
	dataMap, err := k.getMapData()
	if err != nil {
		return nil, err
	}

	return dataMap, nil
}

// Set checks if the value has changed before performing the underlying save call.
func (k *KeyValue) Set(key string, value []byte) error {
	k.Lock()
	defer k.Unlock()

	return k.set(key, value, false)
}

// ForceSet is the same as Set, but does not check if the values are equal first.
func (k *KeyValue) ForceSet(key string, value []byte) error {
	k.Lock()
	defer k.Unlock()

	return k.set(key, value, true)
}

func (k *KeyValue) set(key string, value []byte, force bool) error {
	// Grab the data map.
	dataMap, err := k.getMapData()
	if err != nil {
		return err
	}

	if !force {
		// Look up the original value and check if it's the same.
		if ogValue, ok := dataMap[key]; ok && bytes.Equal(ogValue, value) {
			return nil
		}
	}

	// Set the new value.
	dataMap[key] = value

	// Write the ConfigMap.
	return k.client.Set(k.configmapName, dataMap)
}

// Delete removes the given key from the underlying configmap.
func (k *KeyValue) Delete(key string) error {
	k.Lock()
	defer k.Unlock()

	// Grab the data map.
	dataMap, err := k.getMapData()
	if err != nil {
		return err
	}

	// Delete the key/value.
	delete(dataMap, key)

	// Write the ConfigMap.
	return k.client.Set(k.configmapName, dataMap)
}

// Reset removes all the data from the underlying configmap.
func (k *KeyValue) Reset() error {
	k.Lock()
	defer k.Unlock()

	// Reset the internal cache if needed.
	if k.cacheEnabled {
		k.internalCache = map[string][]byte{}
	}

	// Write the ConfigMap with a new blank map.
	return k.client.Set(k.configmapName, map[string][]byte{})
}
