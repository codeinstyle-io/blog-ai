package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"strings"

	"captain-corp/captain/models"
	"captain-corp/captain/repository"
	"captain-corp/captain/storage"
	"captain-corp/captain/system"

	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"

	"gorm.io/gorm"
)

// ServeMedia serves media files from the configured storage provider
func ServeMedia(repositories *repository.Repositories, storageProvider storage.Provider) fiber.Handler {

	return func(c *fiber.Ctx) error {
		// Get path and trim leading slash if present
		path := c.Path()
		if path == "" {
			return c.Status(http.StatusBadRequest).SendString("No path provided")
		}
		// Trim the leading slash as paths are stored without it in the database
		path = strings.TrimPrefix(path, "/media/")

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
		file, err := storageProvider.Get(path)
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

// GenerateFavicons generates favicon files from a media file
func GenerateFavicons(media *models.Media, storage storage.Provider) error {
	err := media.FetchFile(storage)
	if err != nil {
		return fmt.Errorf("failed to fetch logo file: %w", err)
	}

	// Decode the original image
	origImg, _, err := image.Decode(media.File)
	if err != nil {
		return fmt.Errorf("failed to decode logo file: %w", err)
	}

	// Generate favicon.ico (32x32)
	err = uploadResizedImage(origImg, system.FaviconSize, storage, system.FaviconFilename)
	if err != nil {
		return fmt.Errorf("failed to generate favicon.ico: %w", err)
	}

	// Generate apple-touch-icon.png (180x180)
	err = uploadResizedImage(origImg, system.AppleTouchIconSize, storage, system.AppleTouchIconFilename)
	if err != nil {
		return fmt.Errorf("failed to generate apple-touch-icon.png: %w", err)
	}

	// Generate icon.png (300x300)
	err = uploadResizedImage(origImg, system.IconSize, storage, system.FaviconPngFilename)
	if err != nil {
		return fmt.Errorf("failed to generate favicon.png: %w", err)
	}

	return nil
}

func uploadResizedImage(img image.Image, width int, storage storage.Provider, filename string) error {
	x, y := width, width

	resized := resize.Resize(uint(x), uint(y), img, resize.Lanczos3)
	var buf bytes.Buffer
	if err := png.Encode(&buf, resized); err != nil {
		return err
	}

	_, err := storage.Save(filename, bytes.NewReader(buf.Bytes()))

	return err
}
