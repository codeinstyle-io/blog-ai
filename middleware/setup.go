package middleware

import (
	"net/http"

	"codeinstyle.io/captain/repository"
	"github.com/gin-gonic/gin"
)

// RequireSetup checks if there are any users in the system
func RequireSetup(repos *repository.Repositories) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := repos.Users.FindAll()
		if err != nil || len(users) == 0 {
			if c.Request.URL.Path != "/setup" {
				c.Redirect(http.StatusFound, "/setup")
				c.Abort()
				return
			}
		} else if c.Request.URL.Path == "/setup" {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}
