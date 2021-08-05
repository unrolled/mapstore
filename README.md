# MapStore

[![Test Status](https://github.com/unrolled/mapstore/workflows/tests/badge.svg)](https://github.com/unrolled/mapstore/actions)
[![GoDoc](http://godoc.org/github.com/unrolled/mapstore?status.svg)](http://godoc.org/github.com/unrolled/mapstore)
[![Go Report Card](https://goreportcard.com/badge/github.com/unrolled/mapstore)](https://goreportcard.com/report/github.com/unrolled/mapstore)

## Overview
MapStore is a simple API built on top of Kubernetes ConfigMaps that allows you to programmatically read, write and delete key value pairs. This package was created as a solution to store small amounts of data without needing a persistent volume or standalone database. Because this package is backed by a ConfigMap, we must remember that potentially every interaction will send API requests to the Kubernetes API server.

To get started, check out the [examples](examples) page.

### Caveats
Kubernetes can not guarantee exclusive access to a ConfigMap, so we need to be aware of some edge cases. The ideal usage of MapStore is to have a single process accessing the data to ensure no inconsistences. Using a [leader election](https://github.com/operator-framework/operator-lib/blob/main/leader/doc.go) package to protect access is recommended.

## Internal caching
MapStore has the ability to hold the data of the ConfigMap in memory for quick lookups and reducing unnecessary requests to the Kubernetes API. This should only be enabled when you can guarantee no other app or process is accessing the same ConfigMap.
```go
	cacheConfigMapInternally := true
	mapStore, err := mapstore.New("my-test-cm", cacheConfigMapInternally)
```

## Size limitations
Please be aware that ConfigMaps are limited in size. This package has no protective measures in place to ensure you are below the limit.

> A ConfigMap is not designed to hold large chunks of data. The data stored in a ConfigMap cannot exceed 1 MiB. If you need to store settings that are larger than this limit, you may want to consider mounting a volume or use a separate database or file service.

[Kubernetes ConfigMap documentation](https://kubernetes.io/docs/concepts/configuration/configmap/#motivation)

## Environment variables
There are a few environment variables that you can apply to your workload that will effect MapStore:

`MAPSTORE_CLUSTER_CONFIG_PATH` can be set if you are using this package outside of your cluster, but still want to interact with a ConfigMap on the cluster. You can define the path to your cluster config file and it will be used by MapStore when connecting to the cluster. This is also a handy variable when testing locally. By default this value is empty and MapStore uses the `InClusterConfig` for it's connection.

`NAMESPACE` is the namespace name that MapStore will use. If not set, MapStore will attempt to pull the current namespace from `/var/run/secrets/kubernetes.io/serviceaccount/namespace`.
