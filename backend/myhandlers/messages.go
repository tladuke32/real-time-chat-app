package myhandlers

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/tladuke32/real-time-chat-app/models"
	"log"
	"net/http"
	"time"
)

type NewMessageData struct {
	Content  string `json:"content"`
	UserID   int    `json:"userId"`
	Username string `json:"username"`
}

func HandleNewMessageHTTP(w http.ResponseWriter, r *http.Request) {
	var data NewMessageData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := HandleNewMessage(data.Content, data.UserID, data.Username)
	if err != nil {
		log.Printf("Error handling new message: %v", err)
		http.Error(w, "Error processing message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message created successfully"))
}

func HandleNewMessage(content string, userID int, username string) error {
	// Use GORM to insert the new message
	message := models.Message{
		UserID:   uint(userID),
		Username: username,
		Content:  content,
	}

	if err := db.Create(&message).Error; err != nil {
		log.Printf("Error inserting new message: %v", err)
		return err
	}

	log.Printf("Inserted new message with ID %d for user %d", message.ID, userID)

	// Broadcast the new message to all connected clients
	BroadcastNotification(message)

	return nil
}

func BroadcastNotification(msg models.Message) {
	lock.Lock()
	defer lock.Unlock()

	messageBytes, _ := json.Marshal(msg) // Convert the message to JSON

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Error broadcasting message to WebSocket client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func BroadcastNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Content  string `json:"message"` // Content of the message
		UserID   int    `json:"userId"`  // ID of the user sending the message
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v", err)
		return
	}

	message := models.Message{
		Content:  requestBody.Content,
		UserID:   uint(requestBody.UserID),
		Username: requestBody.Username,
	}

	BroadcastNotification(message)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification broadcasted successfully"))
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding message: %v", err)
		return
	}

	// Use GORM to insert the message
	if err := db.Create(&msg).Error; err != nil {
		http.Error(w, "Error inserting message", http.StatusInternalServerError)
		log.Printf("Error inserting message: %v", err)
		return
	}

	msg.CreatedAt = time.Now()

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []models.Message

	// Use GORM to retrieve messages
	if err := db.Order("created_at DESC").Find(&messages).Error; err != nil {
		http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
		log.Printf("Error retrieving messages: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
