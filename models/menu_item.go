package models

import (
	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	Label    string  `gorm:"not null"`
	URL      *string `gorm:"null"` // External or internal URL
	PageID   *uint   `gorm:"null"` // Reference to Page
	Page     *Page   `gorm:"foreignKey:PageID"`
	Position int     `gorm:"not null;default:0"`
}
