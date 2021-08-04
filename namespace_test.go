package mapstore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamespaceEnv(t *testing.T) {
	nsEnv := "mapstore-test-ns"
	testNamespace := "foobar"

	os.Setenv(nsEnv, testNamespace)

	result, err := getNamespaceFromEnv(nsEnv)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Equal(t, testNamespace, result)
}

func TestNamespaceGoodFile(t *testing.T) {
	testNamespace := "foobar"
	os.Unsetenv(namespaceEnv)

	result, err := getNamespaceFromPath("./testdata/namespace")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Equal(t, testNamespace, result)
}

func TestNamespaceBadFile(t *testing.T) {
	os.Unsetenv(namespaceEnv)

	_, err := getNamespaceFromPath("./testdata/nope-does-not-exist")
	assert.Error(t, err)
}
