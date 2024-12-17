package handlers

import (
	"net/http"

	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gofiber/fiber/v2"
)

// ListPages handles the GET /admin/pages route
func (h *AdminHandlers) ListPages(c *fiber.Ctx) error {
	pages, err := h.repos.Pages.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{})
	}

	return c.Render("admin_pages", fiber.Map{
		"pages": pages,
	})
}

// ShowCreatePage handles the GET /admin/pages/new route
func (h *AdminHandlers) ShowCreatePage(c *fiber.Ctx) error {
	return c.Render("admin_create_page", fiber.Map{
		"page": &models.Page{},
	})
}

// CreatePage handles the POST /admin/pages route
func (h *AdminHandlers) CreatePage(c *fiber.Ctx) error {
	var page models.Page
	if err := c.BodyParser(&page); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_create_page", fiber.Map{
			"page":  &page,
			"error": "Invalid form data",
		})
	}

	// Create page
	if err := h.repos.Pages.Create(&page); err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_create_page", fiber.Map{
			"page":  &page,
			"error": "Failed to create page",
		})
	}

	return c.Redirect("/admin/pages")
}

// EditPage handles the GET /admin/pages/:id/edit route
func (h *AdminHandlers) EditPage(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	page, err := h.repos.Pages.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("admin_edit_page", fiber.Map{
		"page": page,
	})
}

// UpdatePage handles the POST /admin/pages/:id route
func (h *AdminHandlers) UpdatePage(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	page, err := h.repos.Pages.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	// Parse form data
	if err := c.BodyParser(page); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_edit_page", fiber.Map{
			"page":  page,
			"error": "Invalid form data",
		})
	}

	// Update page
	if err := h.repos.Pages.Update(page); err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_edit_page", fiber.Map{
			"page":  page,
			"error": "Failed to update page",
		})
	}

	return c.Redirect("/admin/pages")
}

// ConfirmDeletePage handles the GET /admin/pages/:id/delete route
func (h *AdminHandlers) ConfirmDeletePage(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	page, err := h.repos.Pages.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_page", fiber.Map{
		"page": page,
	})
}

// DeletePage handles the POST /admin/pages/:id/delete route
func (h *AdminHandlers) DeletePage(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page ID"})
	}

	page, err := h.repos.Pages.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Page not found"})
	}

	// Check if page is referenced by menu items
	var menuItemCount int64
	err = h.repos.Pages.CountRelatedMenuItems(page.ID, &menuItemCount)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count related menu items"})
	}

	if menuItemCount > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Page is referenced by menu items"})
	}

	// Delete page
	if err := h.repos.Pages.Delete(page); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete page"})
	}

	return c.JSON(fiber.Map{"message": "Page deleted successfully"})
}
