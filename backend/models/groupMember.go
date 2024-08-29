package models

import "time"

// GroupMember represents a user in a group
type GroupMember struct {
	GroupID   uint      `gorm:"primaryKey;autoIncrement:false" json:"group_id"`
	UserID    uint      `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	Group     Group     `gorm:"foreignKey:GroupID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
