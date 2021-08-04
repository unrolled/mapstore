package mapstore

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	namespaceEnv  = "NAMESPACE"
	namespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func getNamespace() (string, error) {
	envVal, envErr := getNamespaceFromEnv(namespaceEnv)
	if envErr == nil {
		return envVal, nil
	}

	return getNamespaceFromPath(namespacePath)
}

func getNamespaceFromEnv(nsEnv string) (string, error) {
	if val, ok := os.LookupEnv(nsEnv); ok {
		val = strings.TrimSpace(val)
		if len(val) > 0 {
			return val, nil
		}
	}

	return "", fmt.Errorf("no value found in env: %s", nsEnv)
}

func getNamespaceFromPath(nsPath string) (string, error) {
	raw, err := ioutil.ReadFile(nsPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(raw)), nil
}
