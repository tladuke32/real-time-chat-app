package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tladuke32/real-time-chat-app/myhandlers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Setup database connection using environment variables
	done := make(chan bool)
	var db *gorm.DB

	go func() {
		dbUser := os.Getenv("MYSQL_USER")
		dbPassword := os.Getenv("MYSQL_PASSWORD")
		dbHost := os.Getenv("MYSQL_HOST")
		dbPort := os.Getenv("MYSQL_PORT")
		dbName := os.Getenv("MYSQL_DATABASE")
		dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?timeout=5s"

		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("MySQL connection error: %v", err)
			done <- false
		} else {
			log.Println("MySQL connection established.")
			done <- true
		}
	}()

	if !<-done {
		log.Fatal("MySQL connection could not be established. Exiting.")
	}

	myhandlers.MigrateDB(db)

	// Setting up the HTTP router
	r := mux.NewRouter()

	// Group routes based on functionality
	setupUserRoutes(r, db)
	setupMessageRoutes(r, db)
	setupGroupRoutes(r, db)
	setupWebSocketRoutes(r, db)

	// CORS middleware configuration to handle cross-origin requests
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http;//localhost", "http://18.191.149.15"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(r)

	// Start the HTTP server on port 8080
	log.Println("Starting server on :8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
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

func setupUserRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.RegisterUser(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.Login(db, w, r)
	}).Methods("POST")

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/{username}", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.GetUserProfile(db, w, r)
	}).Methods("GET")
	userRouter.HandleFunc("/{userId}/update", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.UpdateUserProfile(db, w, r)
	}).Methods("POST")
}
func setupMessageRoutes(r *mux.Router, db *gorm.DB) {
	messageRouter := r.PathPrefix("/messages").Subrouter()

	messageRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.GetMessages(db, w, r)
	}).Methods("GET")

	messageRouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.HandleNewMessageHTTP(db, w, r)
	}).Methods("POST")

	r.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.SendMessage(db, w, r)
	}).Methods("POST")
}

func setupGroupRoutes(r *mux.Router, db *gorm.DB) {
	groupRouter := r.PathPrefix("/groups").Subrouter()

	groupRouter.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.CreateGroup(db, w, r)
	}).Methods("POST")

	groupRouter.HandleFunc("/add_member", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.AddMemberToGroup(db, w, r)
	}).Methods("POST")

	groupRouter.HandleFunc("/send_message", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.SendMessageToGroup(db, w, r)
	}).Methods("POST")
}

// setupWebSocketRoutes sets up WebSocket routes
func setupWebSocketRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		myhandlers.WebSocketHandler(db, w, r)
	}).Methods("GET")
}
