package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string    `gorm:"not null"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Content     string    `gorm:"not null"`
	PublishedAt time.Time `gorm:"not null"`
	Visible     bool      `gorm:"not null"`
	Excerpt     *string   `gorm:"type:text"`
	Tags        []Tag     `gorm:"many2many:post_tags;"`
	AuthorID    uint      `gorm:"not null"`
	Author      *User     `gorm:"foreignKey:AuthorID"`
}
