package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tladuke32/real-time-chat-app/myhandlers"
	"log"
	"net/http"
	"os"
	"sync"
)

var once sync.Once

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Setup database connection using environment variables
	done := make(chan bool)

	once.Do(func() {
		go func() {
			dbUser := os.Getenv("MYSQL_USER")
			dbPassword := os.Getenv("MYSQL_PASSWORD")
			dbHost := os.Getenv("MYSQL_HOST")
			dbPort := os.Getenv("MYSQL_PORT")
			dbName := os.Getenv("MYSQL_DATABASE")
			dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?timeout=5s"
			myhandlers.InitDB(dsn, done)
		}()
	})

	if !<-done {
		log.Fatal("MySQL connection could not be established. Exiting.")
	}

	// Setting up the HTTP router
	r := mux.NewRouter()

	// Register and login routes
	r.HandleFunc("/register", myhandlers.Register).Methods("POST")
	r.HandleFunc("/login", myhandlers.Login).Methods("POST")

	// WebSocket route for handling real-time communication
	r.HandleFunc("/ws", myhandlers.WebSocketHandler).Methods("GET")

	// User profile management routes
	r.HandleFunc("/user/{username}", myhandlers.GetUserProfile).Methods("GET")
	r.HandleFunc("/user/{userId}/update", myhandlers.UpdateUserProfile).Methods("POST")

	// Message sending and retrieval routes
	r.HandleFunc("/send", myhandlers.SendMessage).Methods("POST")
	r.HandleFunc("/messages", myhandlers.GetMessages).Methods("GET")
	r.HandleFunc("/messages", myhandlers.HandleNewMessageHTTP).Methods("POST")

	// Group management routes
	r.HandleFunc("/groups/create", myhandlers.CreateGroup).Methods("POST")
	r.HandleFunc("/groups/add_member", myhandlers.AddMemberToGroup).Methods("POST")
	r.HandleFunc("/groups/send_message", myhandlers.SendMessageToGroup).Methods("POST")

	// Notification handling route
	r.HandleFunc("/notifications", myhandlers.BroadcastNotificationHandler).Methods("POST")

	// CORS middleware configuration to handle cross-origin requests
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Adjust in production to match your deployment
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "application/json"}),
	)(r)

	// Start the HTTP server on port 8080
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
