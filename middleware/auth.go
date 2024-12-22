package middleware

import (
	"github.com/captain-corp/captain/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// abort redirects the user to the login page
// TODO: Pass the next page
func abort(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})
	return c.Redirect("/login?next=" + c.Path())
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
