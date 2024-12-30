package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/utils"

	"github.com/gofiber/fiber/v2"
)

// ListPosts shows all posts for admin
func (h *AdminHandlers) ListPosts(c *fiber.Ctx) error {
	posts, err := h.repos.Posts.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert post times to user timezone
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(loc)
	}

	return c.Render("admin_posts", fiber.Map{
		"title": "Posts",
		"posts": posts,
	})
}

// ShowCreatePost displays the post creation form
func (h *AdminHandlers) ShowCreatePost(c *fiber.Ctx) error {
	return c.Render("admin_create_post", fiber.Map{
		"title": "Create Post",
	})
}

// ListPostsByTag shows all posts for a specific tag
func (h *AdminHandlers) ListPostsByTag(c *fiber.Ctx) error {
	id := c.Params("id")

	tagID, err := utils.ParseUint(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Find the tag
	tag, err := h.repos.Tags.FindByID(tagID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	posts, err := h.repos.Posts.FindByTag(tag.Slug)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert post times to user timezone
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(loc)
	}

	return c.Render("admin_tag_posts", fiber.Map{
		"title": fmt.Sprintf("Posts tagged with '%s'", tag.Name),
		"posts": posts,
		"tag":   tag,
	})
}

// ConfirmDeletePost shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	tagID, err := utils.ParseUint(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	post, err := h.repos.Posts.FindByID(tagID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_post", fiber.Map{
		"title": "Confirm Post deletion",
		"post":  post,
	})
}

// DeletePost handles post deletion
func (h *AdminHandlers) DeletePost(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid post ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid post ID",
			"redirect": "/admin/posts",
		})
	}

	post, err := h.repos.Posts.FindByID(id)
	if err != nil {
		flash.Error(c, "Post not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":    "Post not found",
			"redirect": "/admin/posts",
		})
	}

	// Delete post
	if err := h.repos.Posts.Delete(post); err != nil {
		flash.Error(c, "Failed to delete post")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to delete post",
			"redirect": "/admin/posts",
		})
	}

	flash.Success(c, "Post deleted successfully")
	return c.JSON(fiber.Map{
		"message":  "Post deleted successfully",
		"redirect": "/admin/posts",
	})
}

func (h *AdminHandlers) ShowEditPost(c *fiber.Ctx) error {
	var post *models.Post
	var err error

	id, err := utils.ParseUint(c.Params("id"))

	if err != nil {
		flash.Error(c, "Invalid post ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid post ID",
			"redirect": "/admin/posts",
		})
	}

	post, err = h.repos.Posts.FindByID(id)

	if err != nil {
		flash.Error(c, "Post not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":    "Post not found",
			"redirect": "/admin/posts",
		})
	}

	return c.Render("admin_edit_post", fiber.Map{
		"title": "Edit Post",
		"post":  post,
	})
}
