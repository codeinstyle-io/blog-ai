package handlers

import (
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/middleware"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"github.com/gin-gonic/gin"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(r *gin.Engine, repos *repository.Repositories, cfg *config.Config) {
	publicHandlers := NewPublicHandlers(repos, cfg)

	// Add setup middleware
	r.Use(middleware.RequireSetup(repos))

	// Auth routes (public)
	r.GET("/", publicHandlers.ListPosts)
	r.GET("/posts/:slug", publicHandlers.GetPostBySlug)
	r.GET("/pages/:slug", publicHandlers.GetPageBySlug)
	r.GET("/tags/:slug", publicHandlers.ListPostsByTag)
	r.GET("/generated/css/chroma.css", publicHandlers.GetChromaCSS)
	r.GET("/media/*path", ServeMedia(repos, cfg))
}

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(r *gin.Engine, repos *repository.Repositories, cfg *config.Config) {
	authHandlers := NewAuthHandlers(repos, cfg)

	r.GET("/admin/setup", authHandlers.HandleSetup)
	r.POST("/admin/setup", authHandlers.HandleSetup)

	// Login routes
	r.GET("/login", authHandlers.Login)
	r.POST("/login", authHandlers.PostLogin)

	// Logout route
	r.GET("/logout", authHandlers.Logout)
}

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(r *gin.Engine, repos *repository.Repositories, cfg *config.Config) {

	storage := storage.NewStorage(cfg)
	adminHandlers := NewAdminHandlers(repos, cfg)
	adminMediaHandlers := NewAdminMediaHandlers(repos, cfg, storage)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(repos))

	// Add index route
	admin.GET("/", adminHandlers.Index)

	// Settings routes
	admin.GET("/settings", adminHandlers.ShowSettings)
	admin.POST("/settings", adminHandlers.UpdateSettings)

	// Tag routes
	admin.GET("/tags", adminHandlers.ListTags)
	admin.GET("/tags/create", adminHandlers.ShowCreateTag)
	admin.POST("/tags/create", adminHandlers.CreateTag)
	admin.GET("/tags/:id/edit", adminHandlers.ShowEditTag)
	admin.POST("/tags/:id/edit", adminHandlers.UpdateTag)
	admin.GET("/tags/:id/posts", adminHandlers.ListPostsByTag)
	admin.GET("/tags/:id/delete", adminHandlers.ConfirmDeleteTag)
	admin.DELETE("/tags/:id", adminHandlers.DeleteTag)

	// User routes
	admin.GET("/users", adminHandlers.ListUsers)
	admin.GET("/users/create", adminHandlers.ShowCreateUser)
	admin.POST("/users/create", adminHandlers.CreateUser)
	admin.GET("/users/:id/edit", adminHandlers.ShowEditUser)
	admin.POST("/users/:id/edit", adminHandlers.UpdateUser)
	admin.GET("/users/:id/delete", adminHandlers.ShowDeleteUser)
	admin.DELETE("/users/:id", adminHandlers.DeleteUser)

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

	// Media routes
	admin.GET("/media", adminMediaHandlers.ListMedia)
	admin.GET("/media/upload", adminMediaHandlers.ShowUploadMedia)
	admin.POST("/media/upload", adminMediaHandlers.UploadMedia)
	admin.GET("/media/:id/delete", adminMediaHandlers.ConfirmDeleteMedia)
	admin.DELETE("/media/:id", adminMediaHandlers.DeleteMedia)

	// API routes
	admin.GET("/api/tags", adminHandlers.GetTags)
	admin.GET("/api/media", adminMediaHandlers.GetMediaList)

}
