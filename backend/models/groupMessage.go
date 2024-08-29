package models

import "github.com/jinzhu/gorm"

type GroupMessage struct {
	gorm.Model
	GroupID uint   `gorm:"not null" json:"group_id"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	Content string `gorm:"type:text;not null" json:"content"`
}
