package myhandlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Upgrader for handling WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // For development, allow all origins. Restrict in production.
	},
}

// Chat handles WebSocket connections and echoes messages back to clients
func Chat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		log.Printf("Error while upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error while reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", msg)

		if err := conn.WriteMessage(messageType, msg); err != nil {
			log.Printf("Error while writing message: %v", err)
			break
		}
	}

	log.Println("Client disconnected")
}
