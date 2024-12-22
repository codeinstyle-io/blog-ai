package middleware

import (
	"fmt"
	"io"
	"strings"

	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/repository"
	"github.com/captain-corp/captain/storage"
	"github.com/captain-corp/captain/system"

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
func ServeFavicon(repositories *repository.Repositories, storage storage.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Serve favicon.ico directly
		var err error
		var data []byte

		if c.Path() == "/"+system.FaviconFilename {
			data, err = readFromStorage(repositories, storage, system.FaviconFilename)
		} else if c.Path() == "/"+system.FaviconPngFilename || c.Path() == "/media/"+system.AppleTouchIconFilename {
			filename := strings.Replace(c.Path(), "/media/", "", 1)
			data, err = readFromStorage(repositories, storage, filename)
		} else {
			return c.Next()
		}

		if err != nil {
			// TODO: Log error
			return c.Next()
		} else {
			if strings.Contains(c.Path(), "favicon.ico") {
				c.Set("Content-Type", "image/x-icon")
			} else {
				c.Set("Content-Type", "image/png")
			}

			c.Set("Cache-Control", "public, max-age=21600")

			return c.Send(data)
		}
	}
}

func readFromStorage(repositories *repository.Repositories, storage storage.Provider, filename string) ([]byte, error) {
	f, err := repositories.Media.FindByFilename(filename)
	fmt.Printf("Favicon: %v\n", f)
	fmt.Printf("Name: %v\n", filename)

	if err != nil {
		return nil, err
	}

	file, err := storage.Get(f.Path) // Remove leading slash
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
