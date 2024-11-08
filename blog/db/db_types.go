package db

import "time"

type Post struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Content     string    `gorm:"not null"`
	PublishedAt time.Time `gorm:"not null"`
	Visible     bool      `gorm:"not null"`
	Excerpt     *string   `gorm:"type:text"` // New nullable field
	Tags        []Tag     `gorm:"many2many:post_tags;"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}
