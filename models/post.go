package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title                     string    `gorm:"not null"`
	Slug                      string    `gorm:"uniqueIndex;not null"`
	Content                   string    `gorm:"not null"`
	PublishedAt               time.Time `gorm:"not null"`
	PublishedAtUTC            time.Time `gorm:"not null"`
	PublishedAtTimezone       string    `gorm:"not null;default:'UTC'"`
	PublishedAtTimeZoneOffset int       `gorm:"not null;default:0"`
	Visible                   bool      `gorm:"not null"`
	Excerpt                   *string   `gorm:"type:text"`
	Tags                      []Tag     `gorm:"many2many:post_tags;"`
	AuthorID                  uint      `gorm:"not null" form:"authorId"`
	Author                    *User     `gorm:"foreignKey:AuthorID" form:"author"`
}

// IsScheduled returns true if the post is scheduled for future publication
func (p *Post) IsScheduled(timezone string) bool {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	return p.Visible && p.PublishedAt.After(now)
}

func (p *Post) ToJSON() string {
	tags := make([]string, 0, len(p.Tags))
	for _, tag := range p.Tags {
		tags = append(tags, tag.Name)
	}

	publishedAt := p.PublishedAt.Format(time.RFC3339)

	buff, err := json.Marshal(map[string]interface{}{
		"id":          p.ID,
		"slug":        p.Slug,
		"title":       p.Title,
		"content":     p.Content,
		"publishedAt": publishedAt,
		"visible":     p.Visible,
		"excerpt":     p.Excerpt,
		"tags":        tags,
		"authorId":    p.AuthorID,
	})
	if err != nil {
		return ""
	}

	return string(buff)
}
