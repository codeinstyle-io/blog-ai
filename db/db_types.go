package db

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Content     string    `gorm:"not null"`
	PublishedAt time.Time `gorm:"not null"`
	Visible     bool      `gorm:"not null"`
	Excerpt     *string   `gorm:"type:text"` // New nullable field
	Tags        []Tag     `gorm:"many2many:post_tags;"`
	AuthorID    uint      `gorm:"not null"`
	Author      *User     `gorm:"foreignKey:AuthorID"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

// User represents a user in the system
type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"uniqueIndex"`
	Password     string
	SessionToken *string `gorm:"uniqueIndex"` // Changed to pointer to allow NULL
}
