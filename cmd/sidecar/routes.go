package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/unrolled/mapstore"
)

func healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func api() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/"), "/")
		if len(pathParts) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid path"))
			return
		}

		switch r.Method {
		case http.MethodGet:
			get(w, pathParts[0], pathParts[1])
		case http.MethodPost:
			set(w, r, pathParts[0], pathParts[1])
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("method not allowed"))
			return
		}
	}
}

func get(w http.ResponseWriter, cmName, key string) {
	mapStore, err := mapstore.NewKeyValue(cmName, false)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	val, err := mapStore.Get(key)

	if err == mapstore.ErrKeyValueNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - not found"))
		return
	} else if err != nil {
		writeError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(val)
}

func set(w http.ResponseWriter, r *http.Request, cmName, key string) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		writeError(w, err.Error())
		return
	}

	mapStore, err := mapstore.NewKeyValue(cmName, false)
	if err != nil {
		writeError(w, err.Error())
		return
	}

	if err := mapStore.Set(key, b); err != nil {
		writeError(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func writeError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}
