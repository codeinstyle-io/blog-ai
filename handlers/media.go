package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"strings"

	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"codeinstyle.io/captain/system"
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
func GenerateFavicons(repositories *repository.Repositories, storage storage.Provider, logoID uint) error {

	logo, err := repositories.Media.FindByID(logoID)
	if err != nil {
		return fmt.Errorf("failed to get logo: %w", err)
	}

	// Read the original image
	origFile, err := storage.Get(logo.Path)
	if err != nil {
		return fmt.Errorf("failed to read logo file: %w", err)
	}

	// Decode the original image
	origImg, _, err := image.Decode(origFile)
	if err != nil {
		return fmt.Errorf("failed to decode logo file: %w", err)
	}

	// Generate favicon.ico (32x32)
	err = uploadResizedImage(origImg, system.FaviconSize, system.FaviconSize, storage, system.FaviconFilename)
	if err != nil {
		return fmt.Errorf("failed to generate favicon.ico: %w", err)
	}

	// Generate apple-touch-icon.png (180x180)
	err = uploadResizedImage(origImg, system.AppleTouchIconSize, system.AppleTouchIconSize, storage, system.AppleTouchIconFilename)
	if err != nil {
		return fmt.Errorf("failed to generate apple-touch-icon.png: %w", err)
	}

	// // If the original is SVG, copy it directly
	// if filepath.Ext(logo.Path) == ".svg" {
	// 	origData, err := storage.Get(logo.Path)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to read original SVG: %w", err)
	// 	}
	// 	if _, err := storage.Save(svgFilename, origData); err != nil {
	// 		return fmt.Errorf("failed to save icon.svg: %w", err)
	// 	}
	// } else {
	// 	// Convert to SVG
	// 	icon, err := oksvg.ReadIconStream(origFile)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to convert to SVG: %w", err)
	// 	}

	// 	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	// 	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// 	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	// 	raster := rasterx.NewDasher(w, h, scanner)
	// 	icon.Draw(raster, 1.0)

	// 	var svgBuf bytes.Buffer
	// 	if err := png.Encode(&svgBuf, img); err != nil {
	// 		return fmt.Errorf("failed to encode SVG: %w", err)
	// 	}
	// 	if _, err := storage.Save(svgFilename, bytes.NewReader(svgBuf.Bytes())); err != nil {
	// 		return fmt.Errorf("failed to save icon.svg: %w", err)
	// 	}
	// }

	return nil
}

func uploadResizedImage(img image.Image, width, height int, storage storage.Provider, filename string) error {
	resized := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	var buf bytes.Buffer
	if err := png.Encode(&buf, resized); err != nil {
		return err
	}

	_, err := storage.Save(filename, bytes.NewReader(buf.Bytes()))

	return err
}
