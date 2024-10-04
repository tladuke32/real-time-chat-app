package myhandlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tladuke32/real-time-chat-app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

// Credentials holds the user's login details
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims structure to hold JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login handles user authentication and JWT generation
func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Decode the JSON request body into the Credentials struct
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v", err)
		return
	}

	// Query the database for the user's stored password
	var user models.User
	err := db.Where("username = ?", creds.Username).First(&user).Error
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Printf("Invalid login attempt for username: %s", creds.Username)
		return
	}

	// Compare the provided password with the stored hashed password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Printf("Invalid password attempt for username: %s", creds.Username)
		return
	}

	// Generate a JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error signing JWT for user %s: %v", creds.Username, err)
		return
	}

	// Set the token as an HTTP-only, Secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production", // Secure only in production
		SameSite: http.SameSiteStrictMode,          // Prevents CSRF
	})

	// Return the token in the response body along with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   tokenString,
	})

	log.Printf("User %s logged in successfully", creds.Username)
}
