/*
Package mapstore facilitates saving key value pairs into a Kubernetes ConfigMap.

  package main

  import (
      "fmt"
      "log"

      "github.com/unrolled/mapstore"
  )

  func main() {
      mapStore, err := mapstore.New("my-custom-config-map-name", false)
      if err != nil {
          log.Fatalf("error creating mapstore (possible rbac issue?): %v", err)
      }

      err = mapStore.Set("my-key", []byte("my value lives here"))
      if err != nil {
          log.Fatalf("error setting value: %v", err)
      }

      val, err := mapStore.Get("my-key")
      if err != nil {
          log.Fatalf("error getting value: %v", err)
      }

      fmt.Printf("Value from ConfigMap: %#v\n", val)
  }
*/
package mapstore
