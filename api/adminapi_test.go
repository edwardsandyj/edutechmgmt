package api

import (
	"context"
	"testing"
)

// MockAdminClient is a mock implementation of the Google Admin API client.
type MockAdminClient struct{}

// MockData is mock data for testing.
var MockData = []Item{
	{ID: "user1", Type: "User", Data: "user1@example.com"},
	{ID: "device1", Type: "Device", Data: "Chromebook"},
}

// getDataFromAdminAPIMock is a mock implementation of getDataFromAdminAPI for testing.
func getDataFromAdminAPIMock(ctx context.Context) ([]Item, error) {
	return MockData, nil
}

func TestGetDataFromAdminAPI(t *testing.T) {
	// Replace the original function with the mock implementation for testing
	getDataFromAdminAPI = getDataFromAdminAPIMock
	defer func() { getDataFromAdminAPI = getDataFromAdminAPIOriginal }()

	// Call the function to be tested
	data, err := getDataFromAdminAPI(context.Background())
	if err != nil {
		t.Fatalf("getDataFromAdminAPI failed: %v", err)
	}

	// Validate the results
	expected := MockData
	if len(data) != len(expected) {
		t.Fatalf("Unexpected number of items. Expected: %d, Actual: %d", len(expected), len(data))
	}
	for i, item := range expected {
		if data[i].ID != item.ID || data[i].Type != item.Type || data[i].Data != item.Data {
			t.Errorf("Unexpected data item at index %d. Expected: %+v, Actual: %+v", i, item, data[i])
		}
	}
}

func TestMain(m *testing.M) {
	// Run tests
	retCode := m.Run()

	// Optionally perform cleanup or other actions after all tests are executed

	// Exit with the status code returned by the tests
	os.Exit(retCode)
}
