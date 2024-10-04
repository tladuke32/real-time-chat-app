package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID       string `gorm:"type:uuid;primary_key"` // UUID as primary key
	UserID   uint   `gorm:"not null"`              // Foreign key reference to User
	Username string `gorm:"type:varchar(100);not null"`
	Content  string `gorm:"type:text;not null"`
	gorm.Model
}

// BeforeCreate hook to generate UUIDs before saving a new message
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New().String() // Generate a new UUID
	return
}
