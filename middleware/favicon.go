package middleware

import (
	"fmt"
	"io"
	"strings"

	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"codeinstyle.io/captain/system"
	"github.com/gofiber/fiber/v2"
)

func generateFaviconHTML() string {
	return fmt.Sprintf(`<link rel="icon" href="/%s" sizes="32x32"><link rel="apple-touch-icon" href="/%s"><link rel="icon" href="/%s" sizes="300x300">`,
		system.FaviconFilename,
		system.AppleTouchIconFilename,
		system.FaviconPngFilename,
	)
}

// InjectFavicon middleware injects favicon HTML into templates
func InjectFavicon(repositories *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Accepts("text/html", "application/xhtml+xml") == "" {
			return c.Next()
		}

		settings := c.Locals("settings").(*models.Settings)

		if settings.UseFavicon {
			err := c.Bind(fiber.Map{
				"faviconHTML": generateFaviconHTML(),
			})

			if err != nil {
				fmt.Printf("Error binding favicon HTML into context: %v\n", err)
			}
		}
		return c.Next()
	}
}

// ServeFavicon middleware serves favicon files from storage
func ServeFavicon(storage storage.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Serve favicon.ico directly
		if c.Path() == "/"+system.FaviconFilename {
			file, err := storage.Get(system.FaviconFilename)
			if err != nil {
				return c.Next()
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				return c.Next()
			}

			c.Set("Content-Type", "image/x-icon")
			return c.Send(data)
		}

		// Serve other favicon files
		if c.Path() == "/"+system.FaviconPngFilename || c.Path() == "/media/"+system.AppleTouchIconFilename {
			filename := strings.Replace(c.Path(), "/media/", "", 1)
			file, err := storage.Get(filename) // Remove leading slash
			if err != nil {
				return c.Next()
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				return c.Next()
			}

			c.Set("Content-Type", "image/png")
			return c.Send(data)
		}

		return c.Next()
	}
}
