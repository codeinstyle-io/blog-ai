package handlers

import (
	"net/http"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandlers struct {
	db     *gorm.DB
	config *config.Config
}

func NewAdminHandlers(database *gorm.DB, config *config.Config) *AdminHandlers {
	return &AdminHandlers{
		db:     database,
		config: config,
	}
}

func (h *AdminHandlers) Index(c *gin.Context) {
	var postCount, tagCount, userCount int64
	var recentPosts []db.Post

	// Get counts
	h.db.Model(&db.Post{}).Count(&postCount)
	h.db.Model(&db.Tag{}).Count(&tagCount)
	h.db.Model(&db.User{}).Count(&userCount)

	// Get 5 most recent posts
	h.db.Order("published_at desc").Limit(5).Find(&recentPosts)

	data := gin.H{
		"title":       "Dashboard",
		"postCount":   postCount,
		"tagCount":    tagCount,
		"userCount":   userCount,
		"recentPosts": recentPosts,
	}

	data = h.addCommonData(c, data)

	c.HTML(http.StatusOK, "admin_index.tmpl", data)
}

// handlers/admin.go
func (h *AdminHandlers) SavePreferences(c *gin.Context) {
	var prefs struct {
		Theme string `json:"theme"`
	}

	if err := c.BindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preferences"})
		return
	}

	// Save theme preference in cookie
	c.SetCookie("admin_theme", prefs.Theme, 3600*24*365, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
