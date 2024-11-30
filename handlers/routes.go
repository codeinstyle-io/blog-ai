package handlers

import (
	"net/http"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(r *gin.Engine, database *gorm.DB, cfg *config.Config) {
	publicHandlers := NewPublicHandlers(database, cfg)
	authHandlers := NewAuthHandlers(database) // Add this

	// Auth routes (public)
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})
	r.POST("/login", authHandlers.Login)
	r.GET("/", publicHandlers.ListPosts)
	r.GET("/posts/:slug", publicHandlers.GetPostBySlug)
	r.GET("/pages/:slug", publicHandlers.GetPageBySlug)
	r.GET("/tags/:tag", publicHandlers.ListPostsByTag)
	r.GET("/generated/css/chroma.css", publicHandlers.GetChromaCSS)
}

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(r *gin.Engine, database *gorm.DB, cfg *config.Config) {
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(database))

	adminHandlers := NewAdminHandlers(database, cfg)
	authHandlers := NewAuthHandlers(database) // Add this line

	// Add index route
	admin.GET("/", adminHandlers.Index)

	// Add logout route
	admin.GET("/logout", authHandlers.Logout)

	// Tag routes
	admin.GET("/tags", adminHandlers.ListTags)
	admin.GET("/tags/:id/posts", adminHandlers.ListPostsByTag) // Add this line
	admin.DELETE("/tags/:id", adminHandlers.DeleteTag)
	admin.GET("/tags/create", adminHandlers.ShowCreateTag)
	admin.POST("/tags/create", adminHandlers.CreateTag)

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

	// Admin routes
	admin.GET("/pages", adminHandlers.ListPages)
	admin.GET("/pages/create", adminHandlers.ShowCreatePage)
	admin.POST("/pages/create", adminHandlers.CreatePage)
	admin.GET("/pages/:id/edit", adminHandlers.EditPage)
	admin.POST("/pages/:id", adminHandlers.UpdatePage)
	admin.DELETE("/pages/:id", adminHandlers.DeletePage)

	// Menu routes
	admin.GET("/menus", adminHandlers.ListMenuItems)
	admin.GET("/menus/create", adminHandlers.ShowCreateMenuItem)
	admin.POST("/menus/create", adminHandlers.CreateMenuItem)
	admin.GET("/menus/:id/edit", adminHandlers.EditMenuItem)
	admin.POST("/menus/:id", adminHandlers.UpdateMenuItem)
	admin.POST("/menus/:id/move/:direction", adminHandlers.MoveMenuItem)
	admin.GET("/menus/:id/delete", adminHandlers.ConfirmDeleteMenuItem)
	admin.POST("/menus/:id/delete", adminHandlers.DeleteMenuItem)

	// Preferences route
	admin.POST("/preferences", adminHandlers.SavePreferences)
}
