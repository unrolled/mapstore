package main

import (
	"fmt"
	"log"

	"github.com/unrolled/mapstore"
)

func main() {
	mapStore, err := mapstore.NewKeyValue("my-custom-config-map-name", false)
	if err != nil {
		log.Fatalf("error creating mapstore: %v", err)
	}

	// Setting the custom value.
	err = mapStore.Set("my-key", []byte("my value lives here"))
	if err != nil {
		log.Fatalf("error setting value: %v", err)
	}

	// ...

	// Getting the custom value.
	val, err := mapStore.Get("my-key")
	if err != nil {
		log.Fatalf("error getting value: %v", err)
	}

	fmt.Printf("Value from ConfigMap: %#v\n", val)
}
