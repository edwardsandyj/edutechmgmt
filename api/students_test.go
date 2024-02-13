package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EduTechMgmt/main/api"	
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAddStudent(t *testing.T) {
	// Mock Redis client
	mockRedisClient := &MockRedisClient{}
	redisClient = mockRedisClient

	// Create a request with a sample student data
	student := Item{ID: "student1", Data: "John Doe"}
	studentJSON, _ := json.Marshal(student)
	req, err := http.NewRequest("POST", "/students", bytes.NewBuffer(studentJSON))
	assert.NoError(t, err)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the addStudent handler function directly
	addStudent(rr, req)

	// Verify the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
  
}

func TestGetStudentsByName(t *testing.T) {
	// Mock Redis client
	mockRedisClient := &MockRedisClient{}
	redisClient = mockRedisClient

	// Create a request with a sample student name
	req, err := http.NewRequest("GET", "/students/name/JohnDoe", nil)
	assert.NoError(t, err)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Mock Gorilla mux router and set URL variables
	router := mux.NewRouter()
	router.HandleFunc("/students/name/{name}", getStudentsByName).Methods("GET")
	router.ServeHTTP(rr, req.WithContext(context.WithValue(req.Context(), mux.VarsKey, map[string]string{"name": "JohnDoe"})))

	// Verify the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the response body
	var students []Item
	err = json.NewDecoder(rr.Body).Decode(&students)
	assert.NoError(t, err)

}


