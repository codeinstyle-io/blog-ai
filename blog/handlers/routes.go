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
	authHandlers := NewAuthHandlers(database) // Add this

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// Auth routes (public)
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})
	r.POST("/login", authHandlers.Login)

	// Other existing routes...
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
	authHandlers := NewAuthHandlers(database) // Add this line

	// Add index route
	admin.GET("/", adminHandlers.Index)

	// Add logout route
	admin.GET("/logout", authHandlers.Logout)

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

	// Preferences route
	admin.POST("/preferences", adminHandlers.SavePreferences)
}
