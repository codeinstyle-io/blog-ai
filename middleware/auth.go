package middleware

import (
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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
func AuthRequired(repos *repository.Repositories, sessionStore *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := sessionStore.Get(c)

		if err != nil {
			return abort(c)
		}

		loggedIn, _ := session.Get("loggedIn").(bool)

		if !loggedIn {
			return abort(c)
		}

		return c.Next()
	}
}
