package main_test

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

func TestAddDevice(t *testing.T) {
	// Mock Redis client
	mockRedisClient := &MockRedisClient{}
	redisClient = mockRedisClient

	// Create a request with a sample device data
	device := Item{ID: "device1", Data: "12345"}
	deviceJSON, _ := json.Marshal(device)
	req, err := http.NewRequest("POST", "/devices", bytes.NewBuffer(deviceJSON))
	assert.NoError(t, err)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the addDevice handler function directly
	addDevice(rr, req)

	// Verify the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify that the device was added to Redis

}

func TestGetDevicesBySerialNumber(t *testing.T) {
	// Mock Redis client
	mockRedisClient := &MockRedisClient{}
	redisClient = mockRedisClient

	// Create a request with a sample serial number
	req, err := http.NewRequest("GET", "/devices/serialnumber/12345", nil)
	assert.NoError(t, err)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Mock Gorilla mux router and set URL variables
	router := mux.NewRouter()
	router.HandleFunc("/devices/serialnumber/{serialnumber}", getDevicesBySerialNumber).Methods("GET")
	router.ServeHTTP(rr, req.WithContext(context.WithValue(req.Context(), mux.VarsKey, map[string]string{"serialnumber": "12345"})))

	// Verify the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the response body
	var devices []Item
	err = json.NewDecoder(rr.Body).Decode(&devices)
	assert.NoError(t, err)

}
