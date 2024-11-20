package middleware

import (
	"net/http"

	"codeinstyle.io/blog/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRequired(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session")
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		user, err := db.GetUserByToken(database, token)
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
