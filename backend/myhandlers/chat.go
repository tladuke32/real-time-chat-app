package myhandlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

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
	defer conn.Close()

	// Register client
	lock.Lock()
	clients[conn] = true
	lock.Unlock()

	log.Println("WebSocket client connected")
	defer log.Println("WebSocket client disconnected")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Broadcast received message to all clients
		broadcastMessage(messageType, message)
	}
}

// broadcastMessage sends messages to all connected clients
func broadcastMessage(messageType int, message []byte) {
	lock.Lock()
	defer lock.Unlock()

	for client := range clients {
		if err := client.WriteMessage(messageType, message); err != nil {
			log.Println("Error writing message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Cleanup on client disconnect
func removeClient(conn *websocket.Conn) {
	lock.Lock()
	delete(clients, conn)
	lock.Unlock()
	conn.Close()
}

func BroadcastNotification(msg string) {
	lock.Lock()
	defer lock.Unlock()

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("Error broadcasting notification:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
