package middleware

import (
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"github.com/gofiber/fiber/v2"
)

func abort(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     system.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})
	return c.Redirect("/login")
}

// AuthRequired ensures that a user is authenticated
func AuthRequired(repos *repository.Repositories) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(system.CookieName)
		if token == "" {
			return abort(c)
		}

		session, err := repos.Sessions.FindByToken(token)
		if err != nil {
			return abort(c)
		}

		user, err := repos.Users.FindByID(session.UserID)
		if err != nil {
			return abort(c)
		}

		c.Locals("user", user)
		return c.Next()
	}
}
