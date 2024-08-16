package myhandlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

var clients = make(map[*websocket.Conn]bool)
var lock = sync.Mutex{}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}
	defer conn.Close()

	lock.Lock()
	clients[conn] = true
	lock.Unlock()

	log.Println("WebSocket client connected")
	defer func() {
		lock.Lock()
		delete(clients, conn)
		lock.Unlock()
		log.Println("WebSocket client disconnected")
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		broadcastMessage(messageType, message)
	}
}

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
