package myhandlers

import (
	"encoding/json"
	"github.com/tladuke32/real-time-chat-app/models"
	"gorm.io/gorm"
	"net/http"
)

// RegisterUser handles user registration
func RegisterUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode the JSON request body into a User struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password before saving (handled by BeforeSave hook in User model)
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	// You can return a success message or generate a JWT token for login (optional)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
