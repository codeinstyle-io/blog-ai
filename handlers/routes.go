package handlers

import (
	"net/http"

	"codeinstyle.io/captain/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(r *gin.Engine, database *gorm.DB) {
	postHandlers := NewPostHandlers(database)

	authHandlers := NewAuthHandlers(database) // Add this

	// Auth routes (public)
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})
	r.POST("/login", authHandlers.Login)
	r.GET("/", postHandlers.ListPosts)
	r.GET("/posts/:slug", postHandlers.GetPostBySlug)
	r.GET("/tags/:tag", postHandlers.ListPostsByTag)
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
