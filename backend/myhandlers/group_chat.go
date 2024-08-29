package myhandlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tladuke32/real-time-chat-app/models"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use GORM to create the group
	if err := db.Create(&group).Error; err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		log.Printf("Error creating group: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func AddMemberToGroup(w http.ResponseWriter, r *http.Request) {
	var member models.GroupMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use GORM to add the member to the group
	if err := db.Create(&member).Error; err != nil {
		http.Error(w, "Failed to add member to group", http.StatusInternalServerError)
		log.Printf("Error adding member to group: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(member)
}

func SendMessageToGroup(w http.ResponseWriter, r *http.Request) {
	var message models.GroupMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use GORM to send the message to the group
	if err := db.Create(&message).Error; err != nil {
		http.Error(w, "Failed to send message to group", http.StatusInternalServerError)
		log.Printf("Error sending message to group: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func FetchGroupMessages(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	var messages []models.GroupMessage
	if err := db.Where("group_id = ?", groupID).Find(&messages).Error; err != nil {
		http.Error(w, "Failed to fetch group messages", http.StatusInternalServerError)
		log.Printf("Error fetching group messages: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
