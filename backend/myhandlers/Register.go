package myhandlers

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"github.com/tladuke32/real-time-chat-app/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v", err)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		log.Printf("Error generating hashed password: %v", err)
		return
	}

	// Insert the user into the database
	db := GetDB()
	user := models.User{
		Username: creds.Username,
		Password: string(hashedPassword),
	}
	err = db.Create(&user).Error

	if err != nil {
		if isUniqueViolationError(err) {
			http.Error(w, "Username already taken", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Printf("Error inserting new user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("User %s registered successfully", creds.Username)
}

// isUniqueViolationError checks if the error is a unique constraint violation
func isUniqueViolationError(err error) bool {
	// MySQL error code 1062 indicates a duplicate entry
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return true
	}
	return false
}
