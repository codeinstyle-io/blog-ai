package models

import (
	"encoding/json"

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

func (p *Page) ToJSON() string {
	buff, err := json.Marshal(map[string]interface{}{
		"id":          p.ID,
		"slug":        p.Slug,
		"title":       p.Title,
		"content":     p.Content,
		"visible":     p.Visible,
		"contentType": p.ContentType,
	})
	if err != nil {
		return ""
	}

	return string(buff)
}
