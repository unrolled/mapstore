package mapstore

import (
	"os"
	"testing"
)

// This value is also set in `testdata/namespace` and needs to match.
const testNamespace = "mapstore"

const testNamespaceEnv = "MAPSTORE_NAMESPACE_TEST"
const testClusterConfigEnv = "MAPSTORE_CLUSTER_CONFIG_PATH_TEST"
const testNamespacePath = "./testdata/namespace"

func TestMain(m *testing.M) {
	// Set our expected env keys to something else for testing.
	namespaceEnv = testNamespaceEnv
	namespacePath = testNamespacePath
	clusterConfigPathEnv = testClusterConfigEnv

	os.Exit(m.Run())
}
