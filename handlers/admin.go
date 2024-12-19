package handlers

import (
	"net/http"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/repository"
	"github.com/gofiber/fiber/v2"
)

// AdminHandlers contains handlers for admin routes
type AdminHandlers struct {
	repos  *repository.Repositories
	config *config.Config
}

// NewAdminHandlers creates a new AdminHandlers instance
func NewAdminHandlers(repos *repository.Repositories, cfg *config.Config) *AdminHandlers {
	return &AdminHandlers{
		repos:  repos,
		config: cfg,
	}
}

// Index handles the GET /admin route
func (h *AdminHandlers) Index(c *fiber.Ctx) error {
	posts, err := h.repos.Posts.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{
			"err": err.Error(),
		})
	}
	postCount := int64(len(posts))

	tags, err := h.repos.Tags.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{
			"err": err.Error(),
		})
	}
	tagCount := int64(len(tags))

	users, err := h.repos.Users.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{
			"err": err.Error(),
		})
	}
	userCount := int64(len(users))

	// Get 5 most recent posts
	recentPosts, err := h.repos.Posts.FindRecent(5)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{
			"err": err.Error(),
		})
	}

	data := fiber.Map{
		"title":       "Dashboard",
		"postCount":   postCount,
		"tagCount":    tagCount,
		"userCount":   userCount,
		"recentPosts": recentPosts,
	}

	return c.Render("admin_index", data)
}
