package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ServeMedia serves media files from the configured storage provider
func ServeMedia(repositories *repository.Repositories, cfg *config.Config) fiber.Handler {
	// Initialize storage provider
	var provider storage.Provider
	var err error

	switch cfg.Storage.Provider {
	case "s3":
		provider, err = storage.NewS3Provider(cfg.Storage.S3.Bucket, cfg.Storage.S3.Region, cfg.Storage.S3.Endpoint, cfg.Storage.S3.AccessKey, cfg.Storage.S3.SecretKey)
	default: // "local"
		provider, err = storage.NewLocalProvider(cfg.Storage.LocalPath)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize storage provider: %v", err))
	}

	return func(c *fiber.Ctx) error {
		// Get path and trim leading slash if present
		path := c.Params("path")
		if path == "" {
			return c.Status(http.StatusBadRequest).SendString("No path provided")
		}
		// Trim the leading slash as paths are stored without it in the database
		path = strings.TrimPrefix(path, "/")

		// Query the media from the database
		media, err := repositories.Media.FindByPath(path)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// If the media doesn't exist, return a 404
				return c.Status(http.StatusNotFound).SendString("Media not found")
			}
			return c.Status(http.StatusInternalServerError).SendString("Error retrieving media")
		}

		// Generate ETag based on last modified time and size
		etag := fmt.Sprintf(`"%x-%x"`, media.UpdatedAt.Unix(), media.Size)

		// Check If-None-Match header
		if match := c.Get("If-None-Match"); match != "" {
			if match == etag {
				return c.Status(http.StatusNotModified).SendString("")
			}
		}

		// Get file from storage provider
		file, err := provider.Get(path)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error retrieving media file")
		}
		defer file.Close()

		// Set content type header
		c.Set("Content-Type", media.MimeType)
		c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", path))
		c.Set("ETag", etag)
		c.Set("Last-Modified", media.UpdatedAt.Format(http.TimeFormat))
		c.Set("Content-Length", fmt.Sprintf("%d", media.Size))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Cache-Control", "public, max-age=31536000")

		// Stream the file to the response
		if _, err := io.Copy(c.Response().BodyWriter(), file); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error streaming media file")
		}

		return nil
	}
}
