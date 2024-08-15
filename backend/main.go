package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	appHandlers "github.com/tladuke32/real-time-chat-app/myhandlers"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", appHandlers.Register).Methods("POST")
	r.HandleFunc("/login", appHandlers.Login).Methods("POST")
	//	r.HandleFunc("/logout", appHandlers.Logout).Methods("POST")
	r.HandleFunc("user/{username}", appHandlers.GetUserProfile).Methods("GET")
	r.HandleFunc("user/{userId}/update", appHandlers.UpdateUserProfile).Methods("POST")
	r.HandleFunc("/ws", appHandlers.Chat).Methods("GET")
	r.HandleFunc("/send", appHandlers.SendMessage).Methods("POST")
	r.HandleFunc("/messages", appHandlers.GetMessages).Methods("GET")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Allow your frontend origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "application/json"}),
	)(r)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))

}
