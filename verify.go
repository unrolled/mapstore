package mapstore

import "fmt"

// VerifyConnection is a helper function that creates a temporary ConfigMap to ensure cluster connectivity and RBAC settings.
func VerifyConnection(testMapName string) error {
	client, err := getKubeClient()
	if err != nil {
		return err
	}

	key := "test"
	val := "ok"
	testData := map[string][]byte{key: []byte(val)}

	// Set a value.
	if err := client.set(testMapName, testData); err != nil {
		return err
	}

	// Get a value.
	if data, err := client.get(testMapName); err != nil {
		return err
	} else if dataVal, ok := data["test"]; !ok || string(dataVal) != val {
		return fmt.Errorf("data is mismatched")
	}

	// Delete a value.
	return client.delete(testMapName)
}
