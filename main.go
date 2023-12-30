package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func exportItemsToCSV(w http.ResponseWriter, r *http.Request) {
	// Retrieve items from the Redis data store
	ctx := context.Background()
	items, err := getItemsFromRedis(ctx, redisClient)
	if err != nil {
		http.Error(w, "Error retrieving items from Redis", http.StatusInternalServerError)
		return
	}

	// Generate CSV file
	filePath := "exported_data.csv"
	err = ExportItemsToCSV(items, filePath)
	if err != nil {
		http.Error(w, "Error exporting items to CSV", http.StatusInternalServerError)
		return
	}

	// Send the CSV file as a response
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	w.Header().Set("Content-Type", "text/csv")
	http.ServeFile(w, r, filePath)
}

func importItemsFromCSV(w http.ResponseWriter, r *http.Request) {
	// Parse the CSV file from the request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error parsing CSV file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read CSV file
	reader := csv.NewReader(file)
	var items []Item
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
			return
		}

		// Parse CSV record into an Item
		if len(record) >= 3 {
			item := Item{
				ID:   record[0],
				Type: record[1],
				Data: record[2],
			}
			items = append(items, item)
		}
	}

	// Update Redis data store with imported items
	ctx := context.Background()
	updateRedisDataStore(ctx, items, redisClient)

	// Respond with a success message
	w.Write([]byte("Import successful"))
}

func editItem(w http.ResponseWriter, r *http.Request) {
	var editedItem Item
	err := json.NewDecoder(r.Body).Decode(&editedItem)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Update Redis data store with the edited item
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", editedItem.Type, editedItem.ID)
	err = redisClient.HSet(ctx, key, "Data", editedItem.Data).Err()
	if err != nil {
		http.Error(w, "Error updating Redis data store", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Write([]byte("Edit successful"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWebSocket)
	router.HandleFunc("/export", exportItemsToCSV).Methods("GET")
	router.HandleFunc("/import", importItemsFromCSV).Methods("POST")
	router.HandleFunc("/edit", editItem).Methods("POST")

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
