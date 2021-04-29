package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/unrolled/mapstore"
)

type myCustomData struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type customWrapper struct {
	kvStore mapstore.KeyValueInterface
}

func (c *customWrapper) Get(key string) (*myCustomData, error) {
	val, err := c.kvStore.Get(key)
	if err != nil {
		return nil, err
	}

	var result *myCustomData
	json.Unmarshal(val, result)

	return result, err
}

func (c *customWrapper) Set(key string, val *myCustomData) error {
	result, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.kvStore.Set(key, result)
}

func main() {
	mapStore, err := mapstore.NewKeyValue("my-custom-config-map-name", false)
	if err != nil {
		log.Fatalf("error creating mapstore: %v", err)
	}

	wrapper := &customWrapper{mapStore}
	initialData := &myCustomData{ID: 1234, UserName: "me", Email: "me@example.com"}

	// Setting the custom value.
	err = wrapper.Set("my-key", initialData)
	if err != nil {
		log.Fatalf("error setting value: %v", err)
	}

	// ...

	// Getting the custom value.
	freshData, err := wrapper.Get("my-key")
	if err != nil {
		log.Fatalf("error getting value: %v", err)
	}

	fmt.Printf("Value from ConfigMap: %#v\n", freshData)
}
