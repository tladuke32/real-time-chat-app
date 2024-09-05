package main

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tladuke32/real-time-chat-app/myhandlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var once sync.Once

func main() {
	// Setup database connection using environment variables
	done := make(chan bool)
	go func() {
		dbUser := os.Getenv("MYSQL_USER")
		dbPassword := os.Getenv("MYSQL_PASSWORD")
		dbHost := os.Getenv("MYSQL_HOST")
		dbPort := os.Getenv("MYSQL_PORT")
		dbName := os.Getenv("MYSQL_DATABASE")
		dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?timeout=5s"
		myhandlers.InitDB(dsn, done)
	}()

	if !<-done {
		log.Fatal("MySQL connection could not be established. Exiting.")
	}

	myhandlers.MigrateDB(myhandlers.GetDB())

	// Setting up the HTTP router
	r := mux.NewRouter()

	// Group routes based on functionality
	setupUserRoutes(r)
	setupMessageRoutes(r)
	setupGroupRoutes(r)
	setupWebSocketRoutes(r)

	// CORS middleware configuration to handle cross-origin requests
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Adjust in production to match your deployment
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "application/json"}),
	)(r)

	// Start the HTTP server on port 8080
	log.Println("Starting server on :8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler,
	}

	// Server run in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	log.Printf("Received %v signal, shutting down server...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}

func setupUserRoutes(r *mux.Router) {
	r.HandleFunc("/register", myhandlers.Register).Methods("POST")
	r.HandleFunc("/login", myhandlers.Login).Methods("POST")
	
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/{username}", myhandlers.GetUserProfile).Methods("GET")
	userRouter.HandleFunc("/{userId}/update", myhandlers.UpdateUserProfile).Methods("POST")
}

func setupMessageRoutes(r *mux.Router) {
	messageRouter := r.PathPrefix("/messages").Subrouter()
	messageRouter.HandleFunc("", myhandlers.GetMessages).Methods("GET")
	messageRouter.HandleFunc("", myhandlers.HandleNewMessageHTTP).Methods("POST")
	r.HandleFunc("/send", myhandlers.SendMessage).Methods("POST")
}

func setupGroupRoutes(r *mux.Router) {
	groupRouter := r.PathPrefix("/groups").Subrouter()
	groupRouter.HandleFunc("/create", myhandlers.CreateGroup).Methods("POST")
	groupRouter.HandleFunc("/add_member", myhandlers.AddMemberToGroup).Methods("POST")
	groupRouter.HandleFunc("/send_message", myhandlers.SendMessageToGroup).Methods("POST")
}

func setupWebSocketRoutes(r *mux.Router) {
	r.HandleFunc("/ws", myhandlers.WebSocketHandler).Methods("GET")
}
