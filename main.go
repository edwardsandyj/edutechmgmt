package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Simulate live updates (replace with actual data retrieval logic)
	for {
		time.Sleep(5 * time.Second)

		message := "Data from the server: " + time.Now().Format("2006-01-02 15:04:05")
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
