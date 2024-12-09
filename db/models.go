package db

import (
	"strings"
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
	Excerpt     *string   `gorm:"type:text"` // New nullable field
	Tags        []Tag     `gorm:"many2many:post_tags;"`
	AuthorID    uint      `gorm:"not null"`
	Author      *User     `gorm:"foreignKey:AuthorID"`
}

type Tag struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null"`
	Slug string `gorm:"uniqueIndex;not null"`
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

// User represents a user in the system
type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"uniqueIndex"`
	Password     string
	SessionToken *string `gorm:"uniqueIndex"` // Changed to pointer to allow NULL
}

type Page struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Slug        string `gorm:"uniqueIndex;not null"`
	Content     string `gorm:"not null"`
	ContentType string `gorm:"not null;default:'markdown'"` // 'markdown' or 'html'
	Visible     bool
}

type MenuItem struct {
	gorm.Model
	Label    string  `gorm:"not null"`
	URL      *string `gorm:"null"` // External or internal URL
	PageID   *uint   `gorm:"null"` // Reference to Page
	Page     *Page   `gorm:"foreignKey:PageID"`
	Position int     `gorm:"not null;default:0"`
}

// Settings represents the site configuration
type Settings struct {
	gorm.Model
	Title        string
	Subtitle     string
	Timezone     string
	ChromaStyle  string
	Theme        string // Add theme field
	PostsPerPage int
}

// Media represents a media file in the system
type Media struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Path        string `gorm:"not null;unique"`
	MimeType    string `gorm:"not null"`
	Size        int64  `gorm:"not null"`
	Description string
}
