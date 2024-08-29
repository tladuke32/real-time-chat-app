package myhandlers

import (
	"github.com/tladuke32/real-time-chat-app/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// InitDB initializes the database connection using GORM
func InitDB(dsn string, done chan<- bool) {
	var err error
	retryCount := 10
	retryDelay := 10 * time.Second

	for retries := retryCount; retries > 0; retries-- {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil && sqlDB.Ping() == nil {
				log.Println("Successfully connected to MySQL")
				done <- true
				return
			}

			if err != nil {
				log.Printf("MySQL connection failed: %v. Retrying in %v... (%d retries left)", err, retryDelay, retries)
			} else {
				log.Printf("MySQL connection failed during Ping: %v. Retrying in %v... (%d retries left)", sqlDB.Ping(), retryDelay, retries)
			}
		} else {
			log.Printf("GORM connection initialization failed: %v. Retrying in %v... (%d retries left)", err, retryDelay, retries)
		}

		time.Sleep(retryDelay)
	}

	if err != nil {
		log.Printf("Failed to connect to MySQL after %d attempts: %v", retryCount, err)
		done <- false
	}
}

// GetDB returns the GORM database connection
func GetDB() *gorm.DB {
	return db
}

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&models.Group{}, &models.GroupMember{}, &models.GroupMessage{}, &models.Message{}, &models.User{})
	if err != nil {
		return
	}

}
