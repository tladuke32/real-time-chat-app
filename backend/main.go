package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tladuke32/real-time-chat-app/myhandlers"
	"log"
	"net/http"
)

func main() {
	// Database setup, consider moving this into a separate setup function or package for better abstraction
	myhandlers.SetupDB("root:secret@tcp(mysql:3306)/chat_app?timeout=5s")

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

	// Group management routes
	r.HandleFunc("/groups/create", myhandlers.CreateGroup).Methods("POST")
	r.HandleFunc("/groups/add_member", myhandlers.AddMemberToGroup).Methods("POST")
	r.HandleFunc("/groups/send_message", myhandlers.SendMessageToGroup).Methods("POST")

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
