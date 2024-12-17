package models

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Slug        string `gorm:"uniqueIndex;not null"`
	Content     string `gorm:"not null"`
	ContentType string `gorm:"not null;default:'markdown'"` // 'markdown' or 'html'
	Visible     bool
}
