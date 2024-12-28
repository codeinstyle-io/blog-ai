package handlers

import (
	"net/http"

	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/utils"

	"github.com/gofiber/fiber/v2"
)

// ListTags handles the GET /admin/tags route
func (h *AdminHandlers) ListTags(c *fiber.Ctx) error {
	tags, err := h.repos.Tags.FindPostsAndCount()
	if err != nil {
		flash.Error(c, "Failed to load tags")
		return c.Status(http.StatusInternalServerError).Render("admin_tags", fiber.Map{})
	}

	return c.Render("admin_tags", fiber.Map{
		"title": "Tags",
		"tags":  tags,
	})
}

// DeleteTag handles tag deletion
func (h *AdminHandlers) DeleteTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid tag ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid tag ID",
			"redirect": "/admin/tags",
		})
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		flash.Error(c, "Tag not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":    "Tag not found",
			"redirect": "/admin/tags",
		})
	}

	// Check if tag has any posts
	count, err := h.repos.Posts.CountByTag(tag.ID)
	if err != nil {
		flash.Error(c, "Failed to check tag usage")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to check tag usage",
			"redirect": "/admin/tags",
		})
	}

	if count > 0 {
		flash.Error(c, "Cannot delete tag that is still in use")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Cannot delete tag that is still in use",
			"redirect": "/admin/tags",
		})
	}

	// Delete tag
	if err := h.repos.Tags.Delete(tag); err != nil {
		flash.Error(c, "Failed to delete tag")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to delete tag",
			"redirect": "/admin/tags",
		})
	}

	flash.Success(c, "Tag deleted successfully")
	return c.JSON(fiber.Map{
		"message":  "Tag deleted successfully",
		"redirect": "/admin/tags",
	})
}

// ConfirmDeleteTag shows deletion confirmation page for a tag
func (h *AdminHandlers) ConfirmDeleteTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_tag", fiber.Map{
		"title": "Confirm Tag deletion",
		"tag":   tag,
	})
}

// ShowCreateTag handles the GET /admin/tags/create route
func (h *AdminHandlers) ShowCreateTag(c *fiber.Ctx) error {
	return c.Render("admin_create_tag", fiber.Map{
		"title": "Create Tag",
	})
}

// CreateTag handles the POST /admin/tags/create route
func (h *AdminHandlers) CreateTag(c *fiber.Ctx) error {
	tag := new(models.Tag)

	if err := c.BodyParser(tag); err != nil {
		flash.Error(c, "Invalid form data")
		return c.Status(http.StatusBadRequest).Render("admin_create_tag", fiber.Map{
			"title": "Tags",
			"tag":   &tag,
		})
	}

	// Create tag
	if err := h.repos.Tags.Create(tag); err != nil {
		flash.Error(c, "Failed to create tag")
		return c.Status(http.StatusInternalServerError).Render("admin_create_tag", fiber.Map{
			"title": "Tags",
			"tag":   &tag,
		})
	}

	flash.Success(c, "Tag created successfully")
	return c.Redirect("/admin/tags")
}

// ShowEditTag handles the GET /admin/tags/:id/edit route
func (h *AdminHandlers) ShowEditTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid tag ID")
		return c.Redirect("/admin/tags")
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		flash.Error(c, "Tag not found")
		return c.Redirect("/admin/tags")
	}

	return c.Render("admin_edit_tag", fiber.Map{
		"title": "Edit Tag",
		"tag":   tag,
	})
}

// UpdateTag handles the POST /admin/tags/:id route
func (h *AdminHandlers) UpdateTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	// Parse form data
	if err := c.BodyParser(&tag); err != nil {
		flash.Error(c, "Invalid form data")
		return c.Status(http.StatusBadRequest).Render("admin_edit_tag", fiber.Map{
			"title": "Tags",
			"tag":   tag,
		})
	}

	// Update tag
	if err := h.repos.Tags.Update(tag); err != nil {
		flash.Error(c, "Failed to update tag")
		return c.Status(http.StatusInternalServerError).Render("admin_edit_tag", fiber.Map{
			"title": "Tags",
			"tag":   tag,
		})
	}

	flash.Success(c, "Tag updated successfully")
	return c.Redirect("/admin/tags")
}
