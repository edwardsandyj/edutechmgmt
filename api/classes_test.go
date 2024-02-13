package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAddClass(t *testing.T) {
	// Mock Redis client
	mockRedisClient := &MockRedisClient{}
	redisClient = mockRedisClient

	// Create a request with a sample class data
	class := Item{ID: "class1", Data: "Class A"}
	classJSON, _ := json.Marshal(class)
	req, err := http.NewRequest("POST", "/classes", bytes.NewBuffer(classJSON))
	assert.NoError(t, err)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the addClass handler function directly
	addClass(rr, req)

	// Verify the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify that the class was added to Redis
	// Add additional assertions based on your Redis logic
}

func TestGetClassesByName(t *testing.T) {
	// Mock Redis client
	mockRedisClient := &MockRedisClient{}
	redisClient = mockRedisClient

	// Create a request with a sample class name
	req, err := http.NewRequest("GET", "/classes/name/ClassA", nil)
	assert.NoError(t, err)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Mock Gorilla mux router and set URL variables
	router := mux.NewRouter()
	router.HandleFunc("/classes/name/{name}", getClassesByName).Methods("GET")
	router.ServeHTTP(rr, req.WithContext(context.WithValue(req.Context(), mux.VarsKey, map[string]string{"name": "ClassA"})))

	// Verify the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the response body
	var classes []Item
	err = json.NewDecoder(rr.Body).Decode(&classes)
	assert.NoError(t, err)

	// Add additional assertions based on your Redis logic and expected response
}
