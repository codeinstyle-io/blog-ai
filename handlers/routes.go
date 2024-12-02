package handlers

import (
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(r *gin.Engine, database *gorm.DB, cfg *config.Config) {
	publicHandlers := NewPublicHandlers(database, cfg)
	authHandlers := NewAuthHandlers(database, cfg) // Add this

	// Auth routes (public)
	r.GET("/", publicHandlers.ListPosts)
	r.GET("/posts/:slug", publicHandlers.GetPostBySlug)
	r.GET("/pages/:slug", publicHandlers.GetPageBySlug)
	r.GET("/tags/:tag", publicHandlers.ListPostsByTag)
	r.GET("/generated/css/chroma.css", publicHandlers.GetChromaCSS)

	// Auth routes (private)
	r.GET("/logout", authHandlers.Logout)
	r.GET("/login", authHandlers.Login)
	r.POST("/login", authHandlers.PostLogin)
}

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(r *gin.Engine, database *gorm.DB, cfg *config.Config) {
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(database))

	adminHandlers := NewAdminHandlers(database, cfg)
	authHandlers := NewAuthHandlers(database, cfg) // Add this line

	// Add index route
	admin.GET("/", adminHandlers.Index)

	// Add logout route
	admin.GET("/logout", authHandlers.Logout)

	// Tag routes
	admin.GET("/tags", adminHandlers.ListTags)
	admin.GET("/tags/create", adminHandlers.ShowCreateTag)
	admin.POST("/tags/create", adminHandlers.CreateTag)
	admin.GET("/tags/:id/posts", adminHandlers.ListPostsByTag) // Add this line
	admin.DELETE("/tags/:id", adminHandlers.DeleteTag)

	// User routes
	admin.GET("/users", adminHandlers.ListUsers)

	// Post routes
	admin.GET("/posts", adminHandlers.ListPosts)
	admin.GET("/posts/create", adminHandlers.ShowCreatePost)
	admin.POST("/posts/create", adminHandlers.CreatePost)
	admin.GET("/posts/:id/edit", adminHandlers.EditPost)
	admin.POST("/posts/:id", adminHandlers.UpdatePost)
	admin.GET("/posts/:id/delete", adminHandlers.ConfirmDeletePost)
	admin.DELETE("/posts/:id", adminHandlers.DeletePost)

	// Admin routes
	admin.GET("/pages", adminHandlers.ListPages)
	admin.GET("/pages/create", adminHandlers.ShowCreatePage)
	admin.POST("/pages/create", adminHandlers.CreatePage)
	admin.GET("/pages/:id/edit", adminHandlers.EditPage)
	admin.POST("/pages/:id", adminHandlers.UpdatePage)
	admin.GET("/pages/:id/delete", adminHandlers.ConfirmDeletePage)
	admin.DELETE("/pages/:id", adminHandlers.DeletePage)

	// Menu routes
	admin.GET("/menus", adminHandlers.ListMenuItems)
	admin.GET("/menus/create", adminHandlers.ShowCreateMenuItem)
	admin.POST("/menus/create", adminHandlers.CreateMenuItem)
	admin.GET("/menus/:id/edit", adminHandlers.EditMenuItem)
	admin.POST("/menus/:id", adminHandlers.UpdateMenuItem)
	admin.POST("/menus/:id/move/:direction", adminHandlers.MoveMenuItem)
	admin.GET("/menus/:id/delete", adminHandlers.ConfirmDeleteMenuItem)
	admin.DELETE("/menus/:id", adminHandlers.DeleteMenuItem)

	// Preferences route
	admin.POST("/preferences", adminHandlers.SavePreferences)

	// API routes
	admin.GET("/api/tags", adminHandlers.GetTags)
}
