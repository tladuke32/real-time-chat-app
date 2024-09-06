package myhandlers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tladuke32/real-time-chat-app/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

var lock sync.Mutex

type NewMessageData struct {
	Content  string `json:"content" validate:"required"`
	UserID   int    `json:"userId" validate:"required"`
	Username string `json:"username" validate:"required"`
}

func HandleNewMessageHTTP(d *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var data NewMessageData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Invalid request body"})
		return
	}

	// Validate data
	if err := validate.Struct(data); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Validation failed", Data: err.Error()})
		return
	}

	// Process and broadcast the new message
	if err := HandleNewMessage(data.Content, data.UserID, data.Username); err != nil {
		log.Printf("Error handling new message: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "Error processing message"})
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{Status: http.StatusCreated, Message: "Message created successfully"})
}

func HandleNewMessage(content string, userID int, username string) error {
	// Insert the new message into the database
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
	messageBytes, _ := json.Marshal(msg) // Convert the message to JSON

	clients.Range(func(key, value interface{}) bool {
		client, ok := key.(*websocket.Conn)
		if !ok {
			return true // Continue to next client
		}
		if err := client.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Error broadcasting message to WebSocket client: %v", err)
			client.Close()
			clients.Delete(client)
		}
		return true
	})
}

func SendMessage(d *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Invalid request payload"})
		log.Printf("Error decoding message: %v", err)
		return
	}

	// Validate message data
	if err := validate.Struct(msg); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Validation failed", Data: err.Error()})
		return
	}

	// Insert the message into the database
	msg.CreatedAt = time.Now()
	if err := db.Create(&msg).Error; err != nil {
		respondWithJSON(w, http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "Error inserting message"})
		log.Printf("Error inserting message: %v", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{Status: http.StatusCreated, Message: "Message sent successfully", Data: msg})
}

func GetMessages(d *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var messages []models.Message

	// Retrieve messages from the database
	if err := db.Order("created_at DESC").Find(&messages).Error; err != nil {
		respondWithJSON(w, http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "Error retrieving messages"})
		log.Printf("Error retrieving messages: %v", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Status: http.StatusOK, Message: "Messages retrieved successfully", Data: messages})
}
