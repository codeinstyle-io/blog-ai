package models

import (
	"strings"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null" form:"name"`
	Slug string `gorm:"uniqueIndex;not null" form:"slug"`
}

// BeforeCreate hook to ensure tag has a slug
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.Slug == "" {
		t.Slug = strings.ToLower(strings.ReplaceAll(t.Name, " ", "-"))
	}
	return nil
}

// BeforeUpdate hook to ensure tag has a slug
func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
	if t.Slug == "" {
		t.Slug = strings.ToLower(strings.ReplaceAll(t.Name, " ", "-"))
	}
	return nil
}
