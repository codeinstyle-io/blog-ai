package handlers

import (
	"net/http"

	"codeinstyle.io/blog/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(r *gin.Engine, database *gorm.DB) {
	postHandlers := NewPostHandlers(database)
	skillsHandlers := NewSkillsHandlers()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	r.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.tmpl", gin.H{})
	})
	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", gin.H{})
	})
	r.GET("/posts", postHandlers.ListPosts)
	r.GET("/posts/:slug", postHandlers.GetPostBySlug)
	r.GET("/skills", skillsHandlers.GetSkills)
}

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(r *gin.Engine, database *gorm.DB) {
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(database))

	adminHandlers := NewAdminHandlers(database)

	// Add index route
	admin.GET("/", adminHandlers.Index)

	// Tag routes
	admin.GET("/tags", adminHandlers.ListTags)
	admin.DELETE("/tags/:id", adminHandlers.DeleteTag)

	// User routes
	admin.GET("/users", adminHandlers.ListUsers)

	// Post routes
	admin.GET("/new_post", adminHandlers.ShowCreatePost)
	admin.POST("/new_post", adminHandlers.CreatePost)
	admin.GET("/posts", adminHandlers.ListPosts)
	admin.GET("/posts/:id/delete", adminHandlers.ConfirmDeletePost)
	admin.DELETE("/posts/:id", adminHandlers.DeletePost)
	admin.GET("/posts/:id/edit", adminHandlers.EditPost)
	admin.POST("/posts/:id", adminHandlers.UpdatePost)
	admin.GET("/api/tags", adminHandlers.GetTags)
}
