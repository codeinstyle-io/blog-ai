package middleware

import (
	"codeinstyle.io/captain/repository"
	"github.com/gofiber/fiber/v2"

	"codeinstyle.io/captain/system"
)

// LoadMenuItems loads menu items into the context
func LoadMenuItems(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		menuItems, err := repos.MenuItems.FindAll()
		if err == nil {
			c.Bind(fiber.Map{"menuItems": menuItems})
		}
		return c.Next()
	}
}

func LoadSettings(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		settings, err := repos.Settings.Get()
		if err == nil {
			c.Bind(fiber.Map{"settings": settings})
		}
		return c.Next()
	}
}

func LoadVersion(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Bind(fiber.Map{"version": system.Version})
		return c.Next()
	}
}

func LoadUserData(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := fiber.Map{}

		c.Bind(data)
		return c.Next()
	}
}
