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

// KeyValueInterface defines the required methods to satisfy the KeyValue implementation.
type KeyValueInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

// Verify we meet the requirements for our own internface ;)
var _ KeyValueInterface = &KeyValue{}

// KeyValue is a thread safe key value store backed by a Kubernetes ConfigMap.
type KeyValue struct {
	*sync.RWMutex
	configmapName string
	client        *KubeClient
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
	var cache map[string][]byte
	if cacheInternally {
		if cm, err := kubeClient.getOrCreateConfigMap(cmName); err != nil {
			return nil, err
		} else {
			cache = cm.BinaryData
		}
	}

	return &KeyValue{
		RWMutex:       &sync.RWMutex{},
		configmapName: cmName,
		client:        kubeClient,
		internalCache: cache,
	}, nil
}

func (k *KeyValue) getMapData() (map[string][]byte, error) {
	if k.internalCache != nil {
		return k.internalCache, nil
	}

	data, err := k.client.Get(k.configmapName)

	if statusError, ok := err.(*errors.StatusError); ok && statusError.Status().Reason == v1.StatusReasonNotFound {
		err = nil
	} else if err != nil {
		return nil, err
	}

	if data == nil {
		data = map[string][]byte{}
	}

	return data, nil
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

// Set checks if the value has changed before performing the underlying save call.
func (k *KeyValue) Set(key string, value []byte) error {
	k.Lock()
	defer k.Unlock()

	// Grab the data map.
	dataMap, err := k.getMapData()
	if err != nil {
		return err
	}

	// Look up the original value and check if it's the same.
	if ogValue, ok := dataMap[key]; ok && bytes.Equal(ogValue, value) {
		return nil
	}

	// Set the new value.
	dataMap[key] = value

	// Write the ConfigMap.
	return k.client.Set(k.configmapName, dataMap)
}

// ForceSet is the same as Set, but does not check if the values are equal first.
func (k *KeyValue) ForceSet(key string, value []byte) error {
	k.Lock()
	defer k.Unlock()

	// Grab the data map.
	dataMap, err := k.getMapData()
	if err != nil {
		return err
	}

	// Set the new value.
	dataMap[key] = value

	// Write the ConfigMap.
	return k.client.Set(k.configmapName, dataMap)
}
