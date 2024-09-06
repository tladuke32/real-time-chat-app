package myhandlers

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/tladuke32/real-time-chat-app/models"
)

var validate = validator.New()

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func CreateGroup(d *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	// Validate group data
	if err := validate.Struct(group); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Validation failed", Data: err.Error()})
		return
	}

	// Use GORM to create the group
	if err := db.Create(&group).Error; err != nil {
		log.Printf("Error creating group: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "Failed to create group"})
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{Status: http.StatusCreated, Message: "Group created successfully", Data: group})
}

func AddMemberToGroup(d *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var member models.GroupMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	// Validate member data
	if err := validate.Struct(member); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Validation failed", Data: err.Error()})
		return
	}

	// Use GORM to add the member to the group
	if err := db.Create(&member).Error; err != nil {
		log.Printf("Error adding member to group: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "Failed to add member to group"})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Status: http.StatusOK, Message: "Member added successfully", Data: member})
}

func SendMessageToGroup(d *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var message models.GroupMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	// Validate message data
	if err := validate.Struct(message); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Validation failed", Data: err.Error()})
		return
	}

	// Use GORM to send the message to the group
	if err := db.Create(&message).Error; err != nil {
		log.Printf("Error sending message to group: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "Failed to send message to group"})
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{Status: http.StatusCreated, Message: "Message sent successfully", Data: message})
}

func FetchGroupMessages(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondWithJSON(w, http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "Group ID is required"})
		return
	}

	var messages []models.GroupMessage
	if err := db.Where("group_id = ?", groupID).Find(&messages).Error; err != nil {
		log.Printf("Error fetching group messages: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "Failed to fetch group messages"})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Status: http.StatusOK, Message: "Messages fetched successfully", Data: messages})
}
