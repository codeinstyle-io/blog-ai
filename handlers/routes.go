package handlers

import (
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/flash"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(repos *repository.Repositories, cfg *config.Config, sessionStore *session.Store) *fiber.App {
	publicHandlers := NewPublicHandlers(repos, cfg)
	app := fiber.New()

	// Public routes
	app.Get("/", publicHandlers.ListPosts)
	app.Get("/posts/:slug", publicHandlers.GetPostBySlug)
	app.Get("/pages/:slug", publicHandlers.GetPageBySlug)
	app.Get("/tags/:slug", publicHandlers.ListPostsByTag)

	return app
}

// RegisterDynamicRoutes registers all dynamic routes
func RegisterDynamicRoutes(repos *repository.Repositories, cfg *config.Config) *fiber.App {
	app := fiber.New()

	app.Get("/chroma.css", GetChromaCSS)
	app.Get("/*", ServeMedia(repos, cfg))

	return app
}

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(repos *repository.Repositories, cfg *config.Config, sessionStore *session.Store) *fiber.App {
	app := fiber.New()
	authHandlers := NewAuthHandlers(repos, cfg, sessionStore)

	app.Get("/setup", authHandlers.HandleSetup)
	app.Post("/setup", authHandlers.HandleSetup)

	// Login routes
	app.Get("/login", authHandlers.ShowLogin)
	app.Post("/login", authHandlers.PostLogin)

	// Logout route
	app.Get("/logout", authHandlers.Logout)

	return app
}

// RegisterAdminRoutes registers all admin routes
func RegisterAdminRoutes(repos *repository.Repositories, cfg *config.Config, sessionStore *session.Store) *fiber.App {
	storage := storage.NewStorage(cfg)
	flash.Setup(sessionStore)
	adminHandlers := NewAdminHandlers(repos, cfg)
	adminMediaHandlers := NewAdminMediaHandlers(repos, cfg, storage)

	app := fiber.New()
	admin := app.Group("/admin")

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
	admin.Get("/users/:id/delete", adminHandlers.ConfirmDeleteUser)
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

	admin.Use("/", flash.Middleware())

	// API routes
	api := admin.Group("/api")
	api.Get("/tags", adminHandlers.GetTags)
	api.Get("/media", adminMediaHandlers.GetMediaList)

	return app
}
