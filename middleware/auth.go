package middleware

import (
	"net/http"

	"codeinstyle.io/captain/repository"
	"github.com/gin-gonic/gin"
)

// AuthRequired ensures that a user is authenticated
func AuthRequired(repos *repository.Repositories) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		user, err := repos.Users.FindBySessionToken(token)
		if err != nil {
			c.SetCookie("session", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
