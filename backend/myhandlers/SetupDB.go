package myhandlers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"log"
	"time"
)

var db *sql.DB

// InitDB initializes the database connection pool
func InitDB(dsn string, done chan<- bool) {
	var err error
	retryCount := 10
	retryDelay := 10 * time.Second

	for retries := retryCount; retries > 0; retries-- {
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			log.Println("Successfully connected to MySQL")
			done <- true
			return
		}

		if err != nil {
			log.Printf("MySQL connection failed: %v. Retrying in 5 seconds... (%d retries left)", err, retries)
		} else {
			log.Printf("MySQL connection failed during Ping: %v. Retrying in 5 seconds... (%d retries left)", db.Ping(), retries)
		}

		time.Sleep(retryDelay)
	}

	if err != nil || db.Ping() != nil {
		log.Printf("Failed to connect to MySQL after %d attempts: %v", retryCount, err)
		done <- false
	}
}

// GetDB returns the database connection pool
func GetDB() *sql.DB {
	return db
}
