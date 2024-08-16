package myhandlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// Upgrader configures WebSocket upgrader options
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Implement more secure origin checking for production environments
	},
}

// clients tracks connected WebSocket clients
var clients = make(map[*websocket.Conn]bool)
var lock = sync.Mutex{} // Ensures concurrent access to clients map is managed safely

// WebSocketHandler handles incoming WebSocket connections
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer func() {
		removeClient(conn)
		conn.Close()
	}()

	// Register the new client
	lock.Lock()
	clients[conn] = true
	lock.Unlock()

	log.Println("WebSocket client connected.")
	defer log.Println("WebSocket client disconnected.")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}

		// Broadcast the received message to all connected clients
		broadcastMessage(messageType, message)
	}
}

// broadcastMessage sends the received message to all connected clients
func broadcastMessage(messageType int, message []byte) {
	lock.Lock()
	defer lock.Unlock()

	for client := range clients {
		if err := client.WriteMessage(messageType, message); err != nil {
			log.Println("Failed to send message:", err)
			removeClient(client) // Ensures clean up of failed
		}
	}
}

func removeClient(conn *websocket.Conn) {
	lock.Lock() // Ensure exclusive access to the clients map
	delete(clients, conn)
	lock.Unlock()
	conn.Close() // Close the WebSocket connection
}
