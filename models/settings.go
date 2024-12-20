package models

import (
	"gorm.io/gorm"
)

// Settings represents the site configuration
type Settings struct {
	gorm.Model
	Title        string `gorm:"not null" form:"title"`
	Subtitle     string `gorm:"not null" form:"subtitle"`
	Timezone     string `gorm:"not null" form:"timezone"`
	ChromaStyle  string `gorm:"not null" form:"chroma_style"`
	Theme        string `gorm:"not null" form:"theme"`
	PostsPerPage int    `gorm:"not null" form:"posts_per_page"`
	LogoID       *uint  `gorm:"" form:"logo_id"`
	UseFavicon   bool   `gorm:"not null;default:false" form:"use_favicon"`
}
