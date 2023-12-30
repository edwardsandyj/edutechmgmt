package main

import (
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
		// Retrieve data from Redis (replace with your actual data retrieval logic)
		ctx := context.Background()

		var items []Item

		deviceKeys, err := redisClient.Keys(ctx, "Device:*").Result()
		if err != nil {
			log.Println("Error retrieving device data from Redis:", err)
			return
		}

		for _, key := range deviceKeys {
			data, err := redisClient.HGet(ctx, key, "Location").Result()
			if err != nil {
				log.Println("Error getting device location:", err)
				continue
			}
			items = append(items, Item{ID: key, Type: "Device", Data: data})
		}

		studentKeys, err := redisClient.Keys(ctx, "Student:*").Result()
		if err != nil {
			log.Println("Error retrieving student data from Redis:", err)
			return
		}

		for _, key := range studentKeys {
			data, err := redisClient.HGet(ctx, key, "Name").Result()
			if err != nil {
				log.Println("Error getting student name:", err)
				continue
			}
			items = append(items, Item{ID: key, Type: "Student", Data: data})
		}

		classKeys, err := redisClient.Keys(ctx, "Class:*").Result()
		if err != nil {
			log.Println("Error retrieving class data from Redis:", err)
			return
		}

		for _, key := range classKeys {
			data, err := redisClient.HGet(ctx, key, "Teacher").Result()
			if err != nil {
				log.Println("Error getting class teacher:", err)
				continue
			}
			items = append(items, Item{ID: key, Type: "Class", Data: data})
		}

		// Send data to the client
		err = conn.WriteJSON(items)
		if err != nil {
			log.Println("Error sending message to client:", err)
			return
		}

		// Wait for the next update
		time.Sleep(5 * time.Second)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
