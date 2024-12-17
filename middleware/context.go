package middleware

import (
	"codeinstyle.io/captain/repository"
	"github.com/gin-gonic/gin"
)

// LoadMenuItems loads menu items into the context
func LoadMenuItems(repos *repository.Repositories) gin.HandlerFunc {
	return func(c *gin.Context) {
		menuItems, err := repos.MenuItems.FindAll()
		if err == nil {
			c.Set("menuItems", menuItems)
		}
		c.Next()
	}
}
