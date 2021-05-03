package main

import (
	"fmt"
	"log"

	"github.com/unrolled/mapstore"
)

const (
	configmapName = "my-custom-config-map-name" // This will be the name of the configmap saved to k8s.
	myKey         = "my-key"                    // This is the key used when writing the key-value pair to the configmap.
)

func main() {
	cacheConfigMapInternally := true
	mapStore, err := mapstore.NewKeyValue(configmapName, cacheConfigMapInternally)
	if err != nil {
		log.Fatalf("error creating mapstore: %v", err)
		// If you receive this error, you likely need to give the appropriate RBAC permissions to your pod.
	}

	// Setting a value. The underlying configmap data has to be a byte slice.
	err = mapStore.Set(myKey, []byte("my value lives here"))
	if err != nil {
		log.Fatalf("error setting value: %v", err)
	}

	// Getting a value.
	val, err := mapStore.Get(myKey)
	if err != nil {
		log.Fatalf("error getting value: %v", err)
	}

	fmt.Printf("Value from ConfigMap: %#v\n", val)
}
