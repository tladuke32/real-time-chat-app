package myhandlers

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

type ChatMessage struct {
	ID        string `json:"id"`        // Unique message ID (to avoid duplicates)
	Username  string `json:"username"`  // Username of the sender
	Message   string `json:"message"`   // Message text
	Timestamp int64  `json:"timestamp"` // Timestamp when the message was sent
}

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

var (
	clients        sync.Map
	messageHistory []ChatMessage
	historyMutex   sync.Mutex
	groupClients   = make(map[string]*sync.Map)
	recentMessages sync.Map
)

func WebSocketHandler(d *gorm.DB, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade error: %v", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	defer func() {
		conn.Close()
		log.Println("WebSocket client disconnected")
		clients.Delete(conn)
	}()

	log.Println("WebSocket client connected")
	clients.Store(conn, true)

	sendChatHistory(conn)

	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		if msg.ID == "" {
			msg.ID = uuid.New().String()
		}

		if _, exists := recentMessages.Load(msg.ID); exists {
			log.Printf("Duplicate message detected: %v", msg.ID)
			continue
		}
		recentMessages.Store(msg.ID, true)

		if msg.Timestamp == 0 {
			msg.Timestamp = time.Now().Unix()
		}

		log.Printf("Broadcasting message from %s: %s", msg.Username, msg.Message)

		historyMutex.Lock()
		messageHistory = append(messageHistory, msg)
		historyMutex.Unlock()

		go broadcastMessage(msg) // Ensure this runs only once per message
	}
}

func sendChatHistory(conn *websocket.Conn) {
	historyMutex.Lock()
	defer historyMutex.Unlock()

	for _, msg := range messageHistory {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending chat history: %v", err)
			break
		}
	}
}

func broadcastMessage(message ChatMessage) {
	clients.Range(func(key, value interface{}) bool {
		client, ok := key.(*websocket.Conn)
		if !ok {
			return true
		}
		if err := client.WriteJSON(message); err != nil {
			log.Printf("Error writing message: %v", err)
			client.Close()
			clients.Delete(client)
		}
		return true
	})
}
