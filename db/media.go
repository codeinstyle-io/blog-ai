package db

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

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

// GetMarkdownTag returns the markdown tag for the media
func (m *Media) GetMarkdownTag() string {
	// If it's an image, return image markdown
	if isImage := isImageMimeType(m.MimeType); isImage {
		return fmt.Sprintf("![%s](/media/%s)", m.Name, m.Path)
	}
	// Otherwise return a link
	return fmt.Sprintf("[%s](/media/%s)", m.Name, m.Path)
}

// GetHTMLTag returns the HTML tag for the media
func (m *Media) GetHTMLTag() string {
	// If it's an image, return img tag
	if isImage := isImageMimeType(m.MimeType); isImage {
		return fmt.Sprintf(`<img src="/media/%s" alt="%s">`, m.Path, m.Name)
	}
	// Otherwise return an anchor tag
	return fmt.Sprintf(`<a href="/media/%s">%s</a>`, m.Path, m.Name)
}

// isImageMimeType checks if the mime type is an image
func isImageMimeType(mimeType string) bool {
	// Check if the mime type starts with "image/"
	return strings.HasPrefix(mimeType, "image/")
}
