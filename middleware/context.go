package middleware

import (
	"fmt"
	"strings"

	"github.com/captain-corp/captain/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/captain-corp/captain/system"
)

func IsAdminPath(c *fiber.Ctx) bool {
	return strings.Index(c.Path(), "/admin") == 0
}

// LoadMenuItems loads menu items into the context
func LoadMenuItems(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if IsAdminPath(c) {
			return c.Next()
		}

		menuItems, err := repos.MenuItems.FindAll()
		if err == nil {
			err = c.Bind(fiber.Map{"menuItems": menuItems})
			if err != nil {
				fmt.Printf("Error binding menu items into context: %v\n", err)
			}

		}
		return c.Next()
	}
}

func LoadSettings(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("settings") != nil {
			return c.Next()
		}

		settings, err := repos.Settings.Get()

		if err == nil {
			c.Locals("settings", settings)
			err = c.Bind(fiber.Map{"settings": settings})
			if err != nil {
				fmt.Printf("Error binding settings into context: %v\n", err)
			}

			c.Locals("settings", settings)
		}
		return c.Next()
	}
}

func LoadVersion(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {

		err := c.Bind(fiber.Map{"version": system.Version})

		if err != nil {
			fmt.Printf("Error binding version into context: %v\n", err)
		}
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
			err = c.Bind(fiber.Map{"user": user})

			if err != nil {
				fmt.Printf("Error binding user into context: %v\n", err)
			}
		}
		return c.Next()
	}
}
