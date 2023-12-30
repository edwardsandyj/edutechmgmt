package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/go-redis/redis/v8"
	"context"
)

var (
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	redisClient *redis.Client
)

type Item struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Data string `json:"data,omitempty"`
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// Retrieve data from Google Workspace Admin API
		ctx := context.Background()
		data, err := getDataFromAdminAPI(ctx)
		if err != nil {
			log.Println("Error retrieving data from Google Workspace Admin API:", err)
			return
		}

		// Send data to the client
		err = conn.WriteJSON(data)
		if err != nil {
			log.Println("Error sending message to client:", err)
			return
		}

		// Update Redis data store with relevant key-value pairs
		updateRedisDataStore(ctx, data, redisClient)

		// Wait for the next update
		time.Sleep(5 * time.Minute)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
