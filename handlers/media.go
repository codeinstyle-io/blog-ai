package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServeMedia serves media files from the media directory
func ServeMedia(database *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get path and trim leading slash if present
		path := c.Param("path")
		if path == "" {
			c.String(http.StatusBadRequest, "No path provided")
			return
		}
		// Trim the leading slash as paths are stored without it in the database
		path = strings.TrimPrefix(path, "/")

		// Query the media from the database
		var media db.Media
		if err := database.Where("path = ?", path).First(&media).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// If the media doesn't exist, return a 404
				c.String(http.StatusNotFound, "Media not found")
				return
			}
			c.String(http.StatusInternalServerError, "Error retrieving media")
			return
		}

		// Construct the full path to the media file
		mediaPath := filepath.Join("media", path)
		if _, err := os.Stat(mediaPath); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "Media file not found")
			return
		}

		// Set content type header
		c.Header("Content-Type", media.MimeType)
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", media.Name))

		// Serve the file
		c.File(mediaPath)
	}
}
