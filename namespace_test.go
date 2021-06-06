package mapstore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamespaceEnv(t *testing.T) {
	os.Setenv(namespaceEnv, testNamespace)

	result, err := getNamespace()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Equal(t, testNamespace, result)
}

func TestNamespaceGoodFile(t *testing.T) {
	os.Unsetenv(namespaceEnv)

	result, err := getNamespace()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Equal(t, testNamespace, result)
}

func TestNamespaceBadFile(t *testing.T) {
	// Makes sure env is blank.
	os.Unsetenv(namespaceEnv)
	// Also set config path to something random.
	namespacePath += "no-such-file"

	_, err := getNamespace()
	assert.Error(t, err)
}
