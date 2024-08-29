package myhandlers

import (
	"encoding/json"
	"errors"
	"github.com/tladuke32/real-time-chat-app/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// GetUserProfile retrieves a user's profile by username
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	var user models.User

	// Use GORM to find the user by username
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error fetching user profile: %v", err)
		return
	}

	// Encode the user profile into JSON and send it in the response
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding user profile response: %v", err)
		return
	}
}

// UpdateUserProfile updates a user's profile information
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Use GORM to update the user's profile
	if err := db.Model(&user).Where("id = ?", user.ID).Updates(models.User{Username: user.Username}).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Printf("Error updating user profile: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
