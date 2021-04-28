package mapstore

import (
	"io/ioutil"
	"os"
	"strings"
)

var namespaceEnv = "NAMESPACE"
var namespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

func getNamespace() (string, error) {
	// Use environment value if possible.
	if val, ok := os.LookupEnv(namespaceEnv); ok {
		val = strings.TrimSpace(val)
		if len(val) > 0 {
			return val, nil
		}
	}

	// Otherwise we need to look it up ourselves.
	raw, err := ioutil.ReadFile(namespacePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(raw)), nil
}
