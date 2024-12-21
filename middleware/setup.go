package middleware

import (
	"captain-corp/captain/repository"

	"github.com/gofiber/fiber/v2"
)

// RequireSetup checks if there are any users in the system
func RequireSetup(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		usersCount, err := repos.Users.CountAll()
		if err != nil {
			return err
		}

		if usersCount > 0 {
			return c.Next()
		}

		if c.Path() == "/setup" {
			return c.Next()
		} else {
			return c.Redirect("/setup")
		}
	}
}
