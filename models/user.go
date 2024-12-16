package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string  `gorm:"uniqueIndex"`
	Password     string
	SessionToken *string `gorm:"uniqueIndex"`
}
