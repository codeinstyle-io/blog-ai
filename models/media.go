package models

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// Media represents a media file in the system
type Media struct {
	gorm.Model
	Name        string `gorm:"not null" form:"name"`
	Path        string `gorm:"not null;unique" form:"path"`
	MimeType    string `gorm:"not null" form:"mimeType"`
	Size        int64  `gorm:"not null" form:"size"`
	Description string `gorm:"type:text" form:"description"`
}

// BeforeCreate hook to ensure media has a mime type
func (m *Media) BeforeCreate(_ *gorm.DB) error {
	if m.MimeType == "" {
		ext := filepath.Ext(m.Name)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		m.MimeType = mimeType
	}
	return nil
}

// GetHTMLTag returns the HTML tag for the media
func (m *Media) GetHTMLTag() string {
	// If it's an image, return img tag
	if isImage := strings.HasPrefix(m.MimeType, "image/"); isImage {
		return fmt.Sprintf(`<img src="/media/%s" alt="%s">`, m.Path, m.Name)
	}
	// Otherwise return an anchor tag
	return fmt.Sprintf(`<a href="/media/%s">%s</a>`, m.Path, m.Name)
}

// GetMarkdownTag returns the markdown tag for the media
func (m *Media) GetMarkdownTag() string {
	// If it's an image, return image markdown
	if isImage := strings.HasPrefix(m.MimeType, "image/"); isImage {
		return fmt.Sprintf("![%s](/media/%s)", m.Name, m.Path)
	}
	// Otherwise return a link
	return fmt.Sprintf("[%s](/media/%s)", m.Name, m.Path)
}
