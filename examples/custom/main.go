package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/unrolled/mapstore"
)

type userObject struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type userStore struct {
	kvStore mapstore.SimpleInterface
}

func (u *userStore) Get(key string) (*userObject, error) {
	val, err := u.kvStore.Get(key)
	if err != nil {
		return nil, err
	}

	var result *userObject
	err = json.Unmarshal(val, result)

	return result, err
}

func (u *userStore) Set(key string, val *userObject) error {
	result, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return u.kvStore.Set(key, result)
}

func main() {
	cacheConfigMapInternally := false
	mapStore, err := mapstore.NewKeyValue("my-custom-config-map-name", cacheConfigMapInternally)
	if err != nil {
		log.Fatalf("error creating mapstore: %v", err)
	}

	wrapper := &userStore{mapStore}
	initialData := &userObject{ID: 1234, UserName: "FooBar", Email: "mapstore@unrolled.ca"}

	// Setting the user data.
	err = wrapper.Set("my-key", initialData)
	if err != nil {
		log.Fatalf("error setting value: %v", err)
	}

	// Getting the user data.
	freshData, err := wrapper.Get("my-key")
	if err != nil {
		log.Fatalf("error getting value: %v", err)
	}

	fmt.Printf("Value from ConfigMap: %#v\n", freshData)
}
