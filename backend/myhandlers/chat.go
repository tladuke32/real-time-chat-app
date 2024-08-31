package myhandlers

import (
		"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//allowedOrigins := map[string]bool{
		//	"http://localhost:3000": true, // Example allowed origin
		//}
		//return allowedOrigins[origin]
		return true
	},
}

var clients sync.Map // Using sync.Map for concurrent-safe access

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade error: %v", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	defer func() {
		conn.Close()
		clients.Delete(conn)
		log.Println("WebSocket client disconnected")
	}()

	log.Println("WebSocket client connected")
	clients.Store(conn, true)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Broadcasting message: %s", string(message)) // Log each broadcast
		go broadcastMessage(messageType, message) // Ensure this runs only once per message
	}
}

func broadcastMessage(messageType int, message []byte) {
	clients.Range(func(key, value interface{}) bool {
		client, ok := key.(*websocket.Conn)
		if !ok {
			return true
		}
		if err := client.WriteMessage(messageType, message); err != nil {
			log.Printf("Error writing message: %v", err)
			client.Close()
			clients.Delete(client)
		}
		return true
	})
}
