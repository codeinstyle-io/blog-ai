package middleware

import (
	"codeinstyle.io/captain/repository"
	"github.com/gofiber/fiber/v2"
)

// LoadMenuItems loads menu items into the context
func LoadMenuItems(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		menuItems, err := repos.MenuItems.FindAll()
		if err == nil {
			c.Bind(fiber.Map{"menuItems": menuItems})
		}
		return nil
	}
}
