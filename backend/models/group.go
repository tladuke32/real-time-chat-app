package models

import "github.com/jinzhu/gorm"

// Group represents a chat group
type Group struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null" json:"name"`
}
