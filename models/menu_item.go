package models

import (
	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	Label    string  `gorm:"not null" json:"label" form:"label"`
	URL      *string `gorm:"null" json:"url" form:"url"`       // External or internal URL
	PageID   *uint   `gorm:"null" json:"pageId" form:"pageId"` // Reference to Page
	Page     *Page   `gorm:"foreignKey:PageID" json:"page" form:"page"`
	Position int     `gorm:"not null;default:0" json:"position" form:"position"`
}
