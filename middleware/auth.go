package middleware

import (
	"net/http"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRequired(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session")
		if err != nil || token == "" {
			// Store requested URL and redirect to login
			next := c.Request.URL.String()
			c.Redirect(http.StatusFound, "/login?next="+next)
			c.Abort()
			return
		}

		user, err := db.GetUserByToken(database, token)
		if err != nil || user.SessionToken == nil {
			c.SetCookie("session", "", -1, "/", "", false, true)
			// Store requested URL and redirect to login
			next := c.Request.URL.String()
			c.Redirect(http.StatusFound, "/login?next="+next)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
