package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockRedisClient is a mock implementation of the Redis client for testing.
type MockRedisClient struct{}

// RedisData is mock data for testing.
var RedisData = map[string]string{
	"item1": "data1",
	"item2": "data2",
}

// MockRedisClient implementation of Get method.
func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	if data, ok := RedisData[key]; ok {
		return data, nil
	}
	return "", nil
}

// MockRedisClient implementation of Keys method.
func (m *MockRedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys := make([]string, 0, len(RedisData))
	for key := range RedisData {
		keys = append(keys, key)
	}
	return keys, nil
}

// MockRedisClient implementation of Set method.
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration int64) error {
	RedisData[key] = value.(string)
	return nil
}

// MockRedisClient implementation of Del method.
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(RedisData, key)
	}
	return nil
}

func TestGetAllItems(t *testing.T) {
	redisClient = &MockRedisClient{}

	req, err := http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllItems)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":"item1","data":"data1"},{"id":"item2","data":"data2"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetItem(t *testing.T) {
	redisClient = &MockRedisClient{}

	req, err := http.NewRequest("GET", "/items/item1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getItem)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"item1","data":"data1"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateItem(t *testing.T) {
	redisClient = &MockRedisClient{}

	item := Item{ID: "item3", Data: "data3"}
	itemJSON, _ := json.Marshal(item)

	req, err := http.NewRequest("POST", "/items", bytes.NewBuffer(itemJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createItem)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if RedisData["item3"] != "data3" {
		t.Errorf("handler did not create item in Redis store")
	}
}

func TestDeleteItem(t *testing.T) {
	redisClient = &MockRedisClient{}

	req, err := http.NewRequest("DELETE", "/items/item1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteItem)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	if _, ok := RedisData["item1"]; ok {
		t.Errorf("handler did not delete item from Redis store")
	}
}
