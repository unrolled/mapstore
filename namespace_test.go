package mapstore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	os.Setenv(namespaceEnv, testNamespace)

	result, err := getNamespace()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Equal(t, testNamespace, result)
}

func TestGoodFile(t *testing.T) {
	os.Unsetenv(namespaceEnv)

	result, err := getNamespace()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	assert.Equal(t, testNamespace, result)
}

func TestBadFile(t *testing.T) {
	// Makes sure env is blank.
	os.Unsetenv(namespaceEnv)
	// Also set config path to something random.
	namespacePath = namespacePath + "no-such-file"

	_, err := getNamespace()
	assert.Error(t, err)
}
