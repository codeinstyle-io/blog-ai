package handlers

import (
	"net/http"

	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gofiber/fiber/v2"
)

// tagResponse struct for API responses
type tagResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// ListTags handles the GET /admin/tags route
func (h *AdminHandlers) ListTags(c *fiber.Ctx) error {
	tags, err := h.repos.Tags.FindPostsAndCount()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{})
	}

	return c.Render("admin_tags", fiber.Map{
		"tags": tags,
	})
}

// DeleteTag handles the DELETE /admin/tags/:id route
func (h *AdminHandlers) DeleteTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag ID"})
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Tag not found"})
	}

	// Check if tag has any posts
	count, err := h.repos.Posts.CountByTag(tag.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check tag usage"})
	}

	if count > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot delete tag with associated posts"})
	}

	// Delete tag
	if err := h.repos.Tags.Delete(tag); err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{})
	}

	return c.Redirect("/admin/tags")
}

// ConfirmDeleteTag shows deletion confirmation page for a tag
func (h *AdminHandlers) ConfirmDeleteTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_tag", fiber.Map{
		"tag": tag,
	})
}

// ShowCreateTag handles the GET /admin/tags/create route
func (h *AdminHandlers) ShowCreateTag(c *fiber.Ctx) error {
	return c.Render("admin_create_tag", fiber.Map{
		"tag": &models.Tag{},
	})
}

// CreateTag handles the POST /admin/tags/create route
func (h *AdminHandlers) CreateTag(c *fiber.Ctx) error {
	tag := new(models.Tag)

	if err := c.BodyParser(tag); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_create_tag", fiber.Map{
			"tag":   &tag,
			"error": "Invalid form data",
		})
	}

	// Create tag
	if err := h.repos.Tags.Create(tag); err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_create_tag", fiber.Map{
			"tag":   &tag,
			"error": "Failed to create tag",
		})
	}

	return c.Redirect("/admin/tags")
}

// ShowEditTag handles the GET /admin/tags/:id/edit route
func (h *AdminHandlers) ShowEditTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("admin_edit_tag", fiber.Map{
		"tag": tag,
	})
}

// UpdateTag handles the POST /admin/tags/:id route
func (h *AdminHandlers) UpdateTag(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	tag, err := h.repos.Tags.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	// Parse form data
	if err := c.BodyParser(&tag); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_edit_tag", fiber.Map{
			"tag":   tag,
			"error": "Invalid form data",
		})
	}

	// Update tag
	if err := h.repos.Tags.Update(tag); err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_edit_tag", fiber.Map{
			"tag":   tag,
			"error": "Failed to update tag",
		})
	}

	return c.Redirect("/admin/tags")
}

// GetTags returns a list of tags for API consumption
func (h *AdminHandlers) GetTags(c *fiber.Ctx) error {
	tags, err := h.repos.Tags.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tags"})
	}

	var response []tagResponse
	for _, tag := range tags {
		response = append(response, tagResponse{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}

	return c.JSON(response)
}
