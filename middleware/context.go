package middleware

import (
	"codeinstyle.io/captain/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

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
			c.Locals("settings", settings)
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

func LoadUserData(repos *repository.Repositories, sessionStore *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := sessionStore.Get(c)

		if err != nil {
			return c.Next()
		}

		userID, _ := session.Get("userID").(uint)
		user, err := repos.Users.FindByID(userID)

		if err == nil {
			c.Locals("user", user)
			c.Bind(fiber.Map{"user": user})
		}
		return c.Next()
	}
}
