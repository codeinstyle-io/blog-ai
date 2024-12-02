package middleware

import (
	"net/http"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireSetup redirects to setup if no users exist
func RequireSetup(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip setup and static assets
		if c.Request.URL.Path == "/admin/setup" || c.Request.URL.Path == "/admin/static" {
			c.Next()
			return
		}

		// Check if any users exist
		var count int64
		database.Model(&db.User{}).Count(&count)
		if count == 0 {
			c.Redirect(http.StatusFound, "/admin/setup")
			c.Abort()
			return
		}

		c.Next()
	}
}
