package middleware

import (
	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoadMenuItems(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []db.MenuItem
		if err := database.Preload("Page").Order("position").Find(&items).Error; err == nil {
			c.Set("menuItems", items)
		}
		c.Next()
	}
}
