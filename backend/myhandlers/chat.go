package myhandlers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type WebSocketMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// Upgrader for handling WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Consider more secure checks for production
	},
}

// Clients map to keep track of connected WebSocket clients
var clients = make(map[*websocket.Conn]bool)
var lock = sync.Mutex{} // to handle concurrent access to clients map

// WebSocketHandler manages all WebSocket connections for notifications and chats
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}

	// Register client
	lock.Lock()
	clients[conn] = true
	lock.Unlock()

	log.Println("WebSocket client connected")
	defer func() {
		removeClient(conn)
		log.Println("WebSocket client disconnected")
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// Broadcast received message to all clients
		broadcastMessage(msg)
	}
}

// broadcastMessage sends messages to all connected clients
func broadcastMessage(msg WebSocketMessage) {
	lock.Lock()
	defer lock.Unlock()

	message, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			removeClient(client)
		}
	}
}

// removeClient handles client disconnection and cleanup
func removeClient(conn *websocket.Conn) {
	lock.Lock()
	defer lock.Unlock()
	delete(clients, conn)
	conn.Close()
}
