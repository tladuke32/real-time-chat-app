package models

import "github.com/jinzhu/gorm"

type Message struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	Username string `gorm:"type:varchar(100);not null"`
	Content  string `gorm:"type:text;not null"`
}
