package myhandlers

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type NewMessageData struct {
	Content string `json:"content"`
	UserID  int    `json:"userId"`
}

func HandleNewMessageHTTP(w http.ResponseWriter, r *http.Request) {
	var data NewMessageData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := HandleNewMessage(data.Content, data.UserID)
	if err != nil {
		log.Printf("Error handling new message: %v", err)
		http.Error(w, "Error processing message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Message created successfully"))
}

func HandleNewMessage(content string, userID int) error {
	// Prepare SQL statement for inserting the new message
	query := `INSERT INTO messages (content, user_id) VALUES (?, ?)`
	result, err := db.Exec(query, content, userID)
	if err != nil {
		log.Printf("Error inserting new message: %v", err)
		return err
	}

	// Retrieve the last inserted ID to confirm the message was saved and for logging purposes
	messageID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID after inserting message: %v", err)
		return err
	}

	log.Printf("Inserted new message with ID %d for user %d", messageID, userID)

	// Optional: Broadcast the new message to all connected clients
	message := fmt.Sprintf("New message from user %d: %s", userID, content)
	BroadcastNotification(message)

	return nil
}

// BroadcastNotification sends a notification to all connected WebSocket clients
func BroadcastNotification(msg string) {
	// Ensure thread safety with a mutex
	lock.Lock()
	defer lock.Unlock()

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Printf("Error broadcasting message to WebSocket client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding message: %v", err)
		return
	}

	query := `INSERT INTO messages (user_id, content) VALUES (?, ?)`
	result, err := db.Exec(query, msg.UserID, msg.Content)
	if err != nil {
		http.Error(w, "Error inserting message", http.StatusInternalServerError)
		log.Printf("Error inserting message: %v", err)
		return
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
		log.Printf("Error getting last insert ID: %v", err)
		return
	}

	msg.ID = int(insertedID)
	msg.CreatedAt = time.Now()

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, user_id, content, created_at FROM messages ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
		log.Printf("Error retrieving messages: %v", err)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Content, &msg.CreatedAt); err != nil {
			http.Error(w, "Error scanning message", http.StatusInternalServerError)
			log.Printf("Error scanning message: %v", err)
			return
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error during rows iteration", http.StatusInternalServerError)
		log.Printf("Error during rows iteration: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
