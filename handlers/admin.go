package handlers

import (
	"net/http"

	"github.com/captain-corp/captain/repository"
	"github.com/captain-corp/captain/storage"

	"github.com/gofiber/fiber/v2"
)

// AdminHandlers contains handlers for admin routes
type AdminHandlers struct {
	repos   *repository.Repositories
	storage storage.Provider
}

// NewAdminHandlers creates a new AdminHandlers instance
func NewAdminHandlers(repos *repository.Repositories, storage storage.Provider) *AdminHandlers {
	return &AdminHandlers{
		repos:   repos,
		storage: storage,
	}
}

// Index handles the GET /admin route
func (h *AdminHandlers) Index(c *fiber.Ctx) error {
	posts, err := h.repos.Posts.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}
	postCount := int64(len(posts))

	pages, err := h.repos.Pages.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}
	pageCount := int64(len(pages))

	tags, err := h.repos.Tags.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}
	tagCount := int64(len(tags))

	users, err := h.repos.Users.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}
	userCount := int64(len(users))

	// Get 5 most recent posts
	recentPosts, err := h.repos.Posts.FindRecent(5)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	data := fiber.Map{
		"title":       "Dashboard",
		"postCount":   postCount,
		"pageCount":   pageCount,
		"tagCount":    tagCount,
		"userCount":   userCount,
		"recentPosts": recentPosts,
	}

	return c.Render("admin_index", data)
}
