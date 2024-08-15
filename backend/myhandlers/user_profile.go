package myhandlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// User represents a user's profile
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"` // Depending on your date handling, this might need to be time.Time
}

// GetUserProfile fetches the profile for a user
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	var user User
	err := db.QueryRow("SELECT id, username, created_at FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error fetching user profile: %v", err)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// UpdateUserProfile updates a user's profile information
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("UPDATE users SET username = ? WHERE id = ?", user.Username, user.ID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Printf("Error updating user profile: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
