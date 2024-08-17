package myhandlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Group represents a chat group
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GroupMember represents a user in a group
type GroupMember struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

// GroupMessage represents a message sent to a group
type GroupMessage struct {
	ID        int    `json:"id"`
	GroupID   int    `json:"group_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"` // Depending on your date handling, this might need to be time.Time
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO `groups` (name) VALUES (?)"
	_, err := db.Exec(query, group.Name)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		log.Printf("Error creating group: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func AddMemberToGroup(w http.ResponseWriter, r *http.Request) {
	var member GroupMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO group_members (group_id, user_id) VALUES (?, ?)"
	_, err := db.Exec(query, member.GroupID, member.UserID)
	if err != nil {
		http.Error(w, "Failed to add member to group", http.StatusInternalServerError)
		log.Printf("Error adding member to group: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(member)
}

func SendMessageToGroup(w http.ResponseWriter, r *http.Request) {
	var message GroupMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO group_messages (group_id, user_id, content) VALUES (?, ?, ?)"
	_, err := db.Exec(query, message.GroupID, message.UserID, message.Content)
	if err != nil {
		http.Error(w, "Failed to send message to group", http.StatusInternalServerError)
		log.Printf("Error sending message to group: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// FetchGroupMessages retrieves all messages for a specific group
func FetchGroupMessages(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT group_id, user_id, content, created_at FROM group_messages WHERE group_id = ?", groupID)
	if err != nil {
		http.Error(w, "Failed to fetch group messages", http.StatusInternalServerError)
		log.Printf("Error fetching group messages: %v", err)
		return
	}
	defer rows.Close()

	var messages []GroupMessage
	for rows.Next() {
		var message GroupMessage
		if err := rows.Scan(&message.GroupID, &message.UserID, &message.Content, &message.CreatedAt); err != nil {
			http.Error(w, "Failed to scan group message", http.StatusInternalServerError)
			log.Printf("Error scanning group message: %v", err)
			return
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error during rows iteration", http.StatusInternalServerError)
		log.Printf("Error during rows iteration: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
