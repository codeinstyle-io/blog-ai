package models

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Title       string `gorm:"not null" form:"title"`
	Slug        string `gorm:"uniqueIndex;not null" form:"slug"`
	Content     string `gorm:"not null" form:"content"`
	ContentType string `gorm:"not null;default:'markdown' " form:"contentType"` // 'markdown' or 'html'
	Visible     bool   `gorm:"not null" form:"visible"`
}
