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
	kvStore mapstore.Interface
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

func main_custom() {
	// Because we are not caching the config map, every request will require going out and looking up the config map.
	// But you need to be aware of the limitations (see main README.md for documentation)!
	cacheConfigMapInternally := false

	mapStore, err := mapstore.New("my-custom-config-map-name", cacheConfigMapInternally)
	if err != nil {
		// If you receive this error, you likely need to give the appropriate RBAC permissions to your pod.
		log.Fatalf("error creating mapstore (possible rbac issue?): %v", err)
	}

	wrapper := &userStore{mapStore}
	initialData := &userObject{ID: 1234, UserName: "FooBar", Email: "mapstore@example.com"}

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
