package handlers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"

	"captain-corp/captain/models"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
	"github.com/gofiber/fiber/v2"
)

func GetChromaCSS(c *fiber.Ctx) error {
	// Generate ETag based on the chroma style name
	settings := c.Locals("settings").(*models.Settings)

	var chromaCSS string
	etag := fmt.Sprintf("\"%x\"", md5.Sum([]byte(settings.ChromaStyle)))

	// Check If-None-Match header first
	if match := c.Get("If-None-Match"); match != "" {
		if match == etag {
			c.Status(http.StatusNotModified)
			return nil
		}
	}

	// Generate CSS
	style := styles.Get(settings.ChromaStyle)
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithClasses(true))
	buf := new(bytes.Buffer)
	if err := formatter.WriteCSS(buf, style); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("")
	}
	chromaCSS = buf.String()

	// Set headers
	c.Set("ETag", etag)
	c.Set("Content-Type", "text/css")
	return c.SendString(chromaCSS)
}
