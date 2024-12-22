package handlers

import (
	"net/http"

	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/utils"

	"github.com/gofiber/fiber/v2"
)

// ListMenuItems displays the menu items management page
func (h *AdminHandlers) ListMenuItems(c *fiber.Ctx) error {
	menuItems, err := h.repos.MenuItems.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	lastPosition := len(menuItems)

	return c.Render("admin_menu_items", fiber.Map{
		"menuItems":    menuItems,
		"title":        "Menu Items",
		"lastPosition": lastPosition,
	})
}

// SaveMenuItems saves the menu items
func (h *AdminHandlers) SaveMenuItems(c *fiber.Ctx) error {
	var menuItems []models.MenuItem
	if err := c.BodyParser(&menuItems); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Delete all existing menu items
	if err := h.repos.MenuItems.DeleteAll(); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Create new menu items
	if err := h.repos.MenuItems.CreateAll(menuItems); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	flash.Success(c, "Menu items saved successfully")
	return c.JSON(fiber.Map{"message": "Menu items saved successfully"})
}

// DeleteMenuItem handles menu item deletion
func (h *AdminHandlers) DeleteMenuItem(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid menu item ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid menu item ID",
			"redirect": "/admin/menus",
		})
	}

	menuItem, err := h.repos.MenuItems.FindByID(id)
	if err != nil {
		flash.Error(c, "Menu item not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":    "Menu item not found",
			"redirect": "/admin/menus",
		})
	}

	// Delete menu item
	if err := h.repos.MenuItems.Delete(menuItem); err != nil {
		flash.Error(c, "Failed to delete menu item")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to delete menu item",
			"redirect": "/admin/menus",
		})
	}

	// Update positions of remaining items
	if err := h.repos.MenuItems.UpdatePositions(menuItem.Position); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    err.Error(),
			"redirect": "/admin/menus",
		})
	}

	flash.Success(c, "Menu item deleted successfully")
	return c.JSON(fiber.Map{
		"message":  "Menu item deleted successfully",
		"redirect": "/admin/menus",
	})
}

// ShowCreateMenuItem displays the menu item creation form
func (h *AdminHandlers) ShowCreateMenuItem(c *fiber.Ctx) error {
	pages, err := h.repos.Pages.FindAll()

	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Render("admin_create_menu_item", fiber.Map{
		"title": "Create Menu Item",
		"pages": pages,
	})
}

// CreateMenuItem handles menu item creation
func (h *AdminHandlers) CreateMenuItem(c *fiber.Ctx) error {
	label := c.FormValue("label")
	urlStr := c.FormValue("url")
	pageID := c.FormValue("page_id")

	if label == "" || (urlStr == "" && pageID == "") {
		flash.Error(c, "Label and either URL or Page are required")
		return c.Status(http.StatusBadRequest).Render("admin_create_menu_item", fiber.Map{
			"title": "Menu Items",
			"pages": []models.Page{},
		})
	}

	menuItem := models.MenuItem{
		Label:    label,
		Position: h.repos.MenuItems.GetNextPosition(),
	}

	// Handle either URL or Page reference
	if pageID != "" {
		pid, _ := utils.ParseUint(pageID)
		menuItem.PageID = &pid
	} else if urlStr != "" {
		menuItem.URL = &urlStr
	}

	if err := h.repos.MenuItems.Create(&menuItem); err != nil {
		pages, err := h.repos.Pages.FindAll()

		if err != nil {
			flash.Error(c, "Failed to create menu item")
			return c.Status(http.StatusInternalServerError).Render("admin_create_menu_item", fiber.Map{
				"title": "Menu Items",
				"item":  menuItem,
				"pages": pages,
			})
		}
	}

	flash.Success(c, "Menu item created successfully")
	return c.Redirect("/admin/menus")
}

// MoveMenuItem handles menu item reordering
func (h *AdminHandlers) MoveMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	direction := c.Params("direction")

	menuID, err := utils.ParseUint(id)
	if err != nil {
		flash.Error(c, "Invalid menu item ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu item ID"})
	}

	if err := h.repos.MenuItems.Move(menuID, direction); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	flash.Success(c, "Menu item moved successfully")
	return c.JSON(fiber.Map{"message": "Menu item moved successfully"})
}

// ConfirmDeleteMenuItem shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeleteMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")

	menuID, err := utils.ParseUint(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	menuItem, err := h.repos.MenuItems.FindByID(menuID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_menu_item", fiber.Map{
		"title":    "Confirm Menu Item deletion",
		"menuItem": menuItem,
	})
}

// EditMenuItem shows the menu item edit form
func (h *AdminHandlers) EditMenuItem(c *fiber.Ctx) error {
	var menuItem *models.MenuItem
	id := c.Params("id")
	menuID, err := utils.ParseUint(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	if menuItem, err = h.repos.MenuItems.FindByID(menuID); err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	pages, err := h.repos.Pages.FindAll()

	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Render("admin_edit_menu_item", fiber.Map{
		"title":    "Edit Menu Item",
		"menuItem": menuItem,
		"pages":    pages,
	})
}

// UpdateMenuItem handles menu item updates
func (h *AdminHandlers) UpdateMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")

	menuID, err := utils.ParseUint(id)
	if err != nil {
		flash.Error(c, "Invalid menu item ID")
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	pages, err := h.repos.Pages.FindAll()
	if err != nil {
		flash.Error(c, "Failed to fetch pages")
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	menuItem, err := h.repos.MenuItems.FindByID(menuID)
	if err != nil {
		flash.Error(c, "Menu item not found")
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	menuItem.Label = c.FormValue("label")
	menuItem.URL = nil
	menuItem.PageID = nil

	// Handle either URL or Page reference
	if pageID := c.FormValue("page_id"); pageID != "" {
		pid, _ := utils.ParseUint(pageID)
		menuItem.PageID = &pid
	} else if urlStr := c.FormValue("url"); urlStr != "" {
		menuItem.URL = &urlStr
	}

	if err := h.repos.MenuItems.Update(menuItem); err != nil {
		flash.Error(c, "Failed to update menu item")
		return c.Status(http.StatusInternalServerError).Render("admin_edit_menu_item", fiber.Map{
			"title": "Menu Items",
			"item":  menuItem,
			"pages": pages,
		})
	}

	flash.Success(c, "Menu item updated successfully")
	return c.Redirect("/admin/menus")
}
