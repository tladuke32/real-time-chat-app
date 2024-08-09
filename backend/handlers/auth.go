package handlers

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

func init() {
	var err error

	// Attempt to connect to the MySQL database with retry logic
	retryCount := 5
	for retries := retryCount; retries > 0; retries-- {
		db, err = sql.Open("mysql", "root:secret@tcp(mysql:3306)/chat-app?timeout=5s")

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

	// Ensure the database connection is still alive
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("MySQL connection is not alive: %v", err))
	}
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error encrypting password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", creds.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles user login and JWT generation
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", creds.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
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
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

// Logout handles user logout by clearing the JWT cookie
//func Logout(w http.ResponseWriter, r *http.Request) {
//	http.SetCookie(w, &http.Cookie{
//		Name:    "token",
//		Value:   "",
//		Expires: time.Now(),
//	})
//	w.WriteHeader(http.StatusOK)
//}
