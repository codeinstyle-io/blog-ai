package middleware

import (
	"net/http"

	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"github.com/gin-gonic/gin"
)

func abort(c *gin.Context) {
	c.SetCookie(system.CookieName, "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
}

// AuthRequired ensures that a user is authenticated
func AuthRequired(repos *repository.Repositories) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(system.CookieName)
		if err != nil {
			abort(c)
			return
		}

		session, err := repos.Sessions.FindByToken(token)
		if err != nil {
			abort(c)
			return
		}

		user, err := repos.Users.FindByID(session.UserID)
		if err != nil {
			abort(c)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
