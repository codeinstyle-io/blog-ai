package middleware

import (
	"net/http"

	"codeinstyle.io/captain/repository"
	"github.com/gofiber/fiber/v2"
)

// RequireSetup checks if there are any users in the system
func RequireSetup(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		usersCount, err := repos.Users.CountAll()
		if err != nil || usersCount == 0 {
			if c.Path() != "/admin/setup" {
				return c.Redirect("/admin/setup", http.StatusFound)
			}
		} else if c.Path() == "/admin/setup" {
			return c.Redirect("/", http.StatusFound)
		}
		return c.Next()
	}
}
