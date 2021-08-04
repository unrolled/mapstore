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
	// By allowing the config map to be stored internally, we reduce the amount of lookups required.
	// But you need to be aware of the limitations (see the main README.md for documentation)!
	cacheConfigMapInternally := true

	// Create a new mapstore.
	mapStore, err := mapstore.New(configmapName, cacheConfigMapInternally)
	if err != nil {
		// If you receive this error, you likely need to give the appropriate RBAC permissions to your pod.
		log.Fatalf("error creating mapstore (possible rbac issue?): %v", err)
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
