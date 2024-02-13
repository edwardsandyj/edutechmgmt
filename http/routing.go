package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Item represents the data model
type Item struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

// getAllItems retrieves all items from Redis
func getAllItems(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, "*").Result()
	if err != nil {
		http.Error(w, "Error getting items", http.StatusInternalServerError)
		return
	}

	var items []Item
	for _, key := range keys {
		data, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			http.Error(w, "Error getting item data", http.StatusInternalServerError)
			return
		}

		items = append(items, Item{ID: key, Data: data})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// getItem retrieves a specific item by ID from Redis
func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	ctx := context.Background()
	data, err := redisClient.Get(ctx, id).Result()
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	item := Item{ID: id, Data: data}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// createItem creates a new item and stores it in Redis
func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)

	ctx := context.Background()
	err := redisClient.Set(ctx, item.ID, item.Data, 0).Err()
	if err != nil {
		http.Error(w, "Error creating item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// deleteItem deletes an item by ID from Redis
func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	ctx := context.Background()
	err := redisClient.Del(ctx, id).Err()
	if err != nil {
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetupRouter configures the API endpoints and starts the server
func SetupRouter() {
	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/items", getAllItems).Methods("GET")
	router.HandleFunc("/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/items", createItem).Methods("POST")
	router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
