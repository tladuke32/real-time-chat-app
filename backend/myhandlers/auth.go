package myhandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var db *sql.DB

// SetupDB initializes the database connection
func SetupDB(dsn string) {
	var err error
	retryCount := 5
	for retries := retryCount; retries > 0; retries-- {
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			log.Println("Successfully connected to MySQL")
			break
		}
		log.Printf("MySQL connection failed: %v. Retrying in 5 seconds... (%d retries left)", err, retries)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL after %d attempts: %v", retryCount, err))
	}
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		log.Printf("Error generating hashed password: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", creds.Username, string(hashedPassword))
	if err != nil {
		if isUniqueViolationError(err) {
			http.Error(w, "Username already taken", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error inserting new user: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("User %s registered successfully", creds.Username)
}

// Login handles user authentication and JWT generation
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request payload: %v", err)
		return
	}

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", creds.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			log.Printf("Invalid login attempt for username: %s", creds.Username)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error querying database for user %s: %v", creds.Username, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Printf("Invalid password attempt for username: %s", creds.Username)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error signing JWT for user %s: %v", creds.Username, err)
		return
	}

	// Set the token as an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,                    // Security best practice to prevent JavaScript access
		SameSite: http.SameSiteStrictMode, // Security best practice
	})

	// Also return the token in the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})

	log.Printf("User %s logged in successfully", creds.Username)
}

// isUniqueViolationError checks if the error is a unique constraint violation
func isUniqueViolationError(err error) bool {
	// Example for MySQL: error code 1062 indicates a duplicate entry
	if err != nil && err.Error() != "" && err.Error() == "Error 1062: Duplicate entry" {
		return true
	}
	return false
}
