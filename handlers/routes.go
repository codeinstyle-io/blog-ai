package handlers

import (
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/middleware"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"github.com/gofiber/fiber/v2"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(app *fiber.App, repos *repository.Repositories, cfg *config.Config) {
	publicHandlers := NewPublicHandlers(repos, cfg)

	// Add setup middleware
	app.Use(middleware.RequireSetup(repos))

	// Public routes
	app.Get("/", publicHandlers.ListPosts)
	app.Get("/posts/:slug", publicHandlers.GetPostBySlug)
	app.Get("/pages/:slug", publicHandlers.GetPageBySlug)
	app.Get("/tags/:slug", publicHandlers.ListPostsByTag)
	app.Get("/generated/css/chroma.css", publicHandlers.GetChromaCSS)
	app.Get("/media/*", ServeMedia(repos, cfg))
}

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(app *fiber.App, repos *repository.Repositories, cfg *config.Config) {
	authHandlers := NewAuthHandlers(repos, cfg)

	app.Get("/admin/setup", authHandlers.HandleSetup)
	app.Post("/admin/setup", authHandlers.HandleSetup)

	// Login routes
	app.Get("/login", authHandlers.ShowLogin)
	app.Post("/login", authHandlers.PostLogin)

	// Logout route
	app.Get("/logout", authHandlers.Logout)
}

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(app *fiber.App, repos *repository.Repositories, cfg *config.Config) {
	storage := storage.NewStorage(cfg)
	adminHandlers := NewAdminHandlers(repos, cfg)
	adminMediaHandlers := NewAdminMediaHandlers(repos, cfg, storage)

	admin := app.Group("/admin")
	admin.Use(middleware.AuthRequired(repos))

	// Dashboard
	admin.Get("/", adminHandlers.Index)

	// Posts
	admin.Get("/posts", adminHandlers.ListPosts)
	admin.Get("/posts/create", adminHandlers.ShowCreatePost)
	admin.Post("/posts/create", adminHandlers.CreatePost)
	admin.Get("/posts/:id/edit", adminHandlers.EditPost)
	admin.Post("/posts/:id", adminHandlers.UpdatePost)
	admin.Get("/posts/:id/delete", adminHandlers.ConfirmDeletePost)
	admin.Delete("/posts/:id", adminHandlers.DeletePost)

	// Pages
	admin.Get("/pages", adminHandlers.ListPages)
	admin.Get("/pages/create", adminHandlers.ShowCreatePage)
	admin.Post("/pages/create", adminHandlers.CreatePage)
	admin.Get("/pages/:id/edit", adminHandlers.EditPage)
	admin.Post("/pages/:id", adminHandlers.UpdatePage)
	admin.Get("/pages/:id/delete", adminHandlers.ConfirmDeletePage)
	admin.Delete("/pages/:id", adminHandlers.DeletePage)

	// Tags
	admin.Get("/tags", adminHandlers.ListTags)
	admin.Get("/tags/create", adminHandlers.ShowCreateTag)
	admin.Post("/tags/create", adminHandlers.CreateTag)
	admin.Get("/tags/:id/edit", adminHandlers.ShowEditTag)
	admin.Post("/tags/:id/edit", adminHandlers.UpdateTag)
	admin.Get("/tags/:id/posts", adminHandlers.ListPostsByTag)
	admin.Get("/tags/:id/delete", adminHandlers.ConfirmDeleteTag)
	admin.Delete("/tags/:id", adminHandlers.DeleteTag)

	// Users
	admin.Get("/users", adminHandlers.ListUsers)
	admin.Get("/users/create", adminHandlers.ShowCreateUser)
	admin.Post("/users/create", adminHandlers.CreateUser)
	admin.Get("/users/:id/edit", adminHandlers.ShowEditUser)
	admin.Post("/users/:id/edit", adminHandlers.UpdateUser)
	admin.Get("/users/:id/delete", adminHandlers.ShowDeleteUser)
	admin.Delete("/users/:id", adminHandlers.DeleteUser)

	// Menus
	admin.Get("/menus", adminHandlers.ListMenuItems)
	admin.Get("/menus/create", adminHandlers.ShowCreateMenuItem)
	admin.Post("/menus/create", adminHandlers.CreateMenuItem)
	admin.Get("/menus/:id/edit", adminHandlers.EditMenuItem)
	admin.Post("/menus/:id", adminHandlers.UpdateMenuItem)
	admin.Post("/menus/:id/move/:direction", adminHandlers.MoveMenuItem)
	admin.Get("/menus/:id/delete", adminHandlers.ConfirmDeleteMenuItem)
	admin.Delete("/menus/:id", adminHandlers.DeleteMenuItem)

	// Media
	admin.Get("/media", adminMediaHandlers.ListMedia)
	admin.Get("/media/upload", adminMediaHandlers.ShowUploadMedia)
	admin.Post("/media/upload", adminMediaHandlers.UploadMedia)
	admin.Get("/media/:id/delete", adminMediaHandlers.ConfirmDeleteMedia)
	admin.Delete("/media/:id", adminMediaHandlers.DeleteMedia)

	// Settings
	admin.Get("/settings", adminHandlers.ShowSettings)
	admin.Post("/settings", adminHandlers.UpdateSettings)

	// API routes
	admin.Get("/api/tags", adminHandlers.GetTags)
	admin.Get("/api/media", adminMediaHandlers.GetMediaList)
}
