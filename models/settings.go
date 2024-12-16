package models

import (
	"gorm.io/gorm"
)

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
