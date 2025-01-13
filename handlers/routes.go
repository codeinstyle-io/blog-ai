package handlers

import (
	"github.com/captain-corp/captain/config"
	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/repository"
	"github.com/captain-corp/captain/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// RegisterPublicRoutes registers all public routes
func RegisterPublicRoutes(repos *repository.Repositories, cfg *config.Config) *fiber.App {
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
func RegisterDynamicRoutes(repos *repository.Repositories, storageProvider storage.Provider) *fiber.App {
	app := fiber.New()

	app.Get("/chroma.css", GetChromaCSS)
	app.Get("/*", ServeMedia(repos, storageProvider))

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
func RegisterAdminRoutes(repos *repository.Repositories, storage storage.Provider, sessionStore *session.Store) *fiber.App {

	flash.Setup(sessionStore)
	adminHandlers := NewAdminHandlers(repos, storage)
	adminMediaHandlers := NewAdminMediaHandlers(repos, storage)

	app := fiber.New()
	admin := app.Group("/admin")

	// Dashboard
	admin.Get("/", adminHandlers.Index)

	// Posts
	admin.Get("/posts", adminHandlers.ListPosts)
	admin.Get("/posts/create", adminHandlers.ShowCreatePost)
	admin.Get("/posts/:id/edit", adminHandlers.ShowEditPost)
	admin.Get("/posts/:id/delete", adminHandlers.ConfirmDeletePost)
	admin.Delete("/posts/:id", adminHandlers.DeletePost)

	// Pages
	admin.Get("/pages", adminHandlers.ListPages)
	admin.Get("/pages/create", adminHandlers.ShowCreatePage)
	admin.Get("/pages/:id/edit", adminHandlers.EditPage)
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
	api.Get("/tags", adminHandlers.ApiGetTags)
	api.Get("/media", adminMediaHandlers.ApiGetMediaList)

	// Posts API routes
	api.Post("/posts", adminHandlers.ApiCreatePost)
	api.Put("/posts/:id", adminHandlers.ApiUpdatePost)

	// Pages API routes
	api.Post("/pages", adminHandlers.ApiCreatePage)
	api.Put("/pages/:id", adminHandlers.ApiUpdatePage)

	return app
}
