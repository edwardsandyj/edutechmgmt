package main

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

// MockRedisClient is a mock implementation of the Redis client for testing.
type MockRedisClient struct {
	Data map[string]string
}

// HSet simulates the HSet operation of the Redis client.
func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.StatusCmd {
	if len(values) != 2 {
		return redis.NewStatusCmd("error", "Invalid number of arguments")
	}
	field, ok := values[0].(string)
	if !ok {
		return redis.NewStatusCmd("error", "Invalid field type")
	}
	value, ok := values[1].(string)
	if !ok {
		return redis.NewStatusCmd("error", "Invalid value type")
	}
	m.Data[key+":"+field] = value
	return redis.NewStatusCmd("OK", "")
}

func TestUpdateRedisDataStore(t *testing.T) {
	ctx := context.Background()

	// Initialize the mock Redis client
	mockClient := &MockRedisClient{
		Data: make(map[string]string),
	}

	// Test data
	data := []Item{
		{Type: "user", ID: "1", Data: "user1"},
		{Type: "user", ID: "2", Data: "user2"},
		{Type: "user", ID: "3", Data: "user3"},
	}

	// Call the function to be tested
	updateRedisDataStore(ctx, data, mockClient)

	// Validate the results
	for _, item := range data {
		key := item.Type + ":" + item.ID
		if value, ok := mockClient.Data[key+":Data"]; !ok {
			t.Errorf("Key %s not found in Redis data store", key)
		} else if value != item.Data {
			t.Errorf("Unexpected value for key %s in Redis data store. Expected: %s, Actual: %s", key, item.Data, value)
		}
	}
}
