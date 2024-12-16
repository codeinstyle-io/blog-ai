package handlers

import (
	"net/http"
	"strconv"

	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListMenuItems displays the menu items management page
func (h *AdminHandlers) ListMenuItems(c *gin.Context) {
	menuItems, err := h.menuRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_menu.tmpl", h.addCommonData(c, gin.H{
		"menuItems": menuItems,
	}))
}

// SaveMenuItems saves the menu items
func (h *AdminHandlers) SaveMenuItems(c *gin.Context) {
	var menuItems []models.MenuItem
	if err := c.ShouldBindJSON(&menuItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete all existing menu items
	if err := h.menuRepo.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create new menu items
	if err := h.menuRepo.CreateAll(menuItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu items saved successfully"})
}

// DeleteMenuItem deletes a menu item
func (h *AdminHandlers) DeleteMenuItem(c *gin.Context) {
	menuItemID, err := utils.ParseUint(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	menuItem, err := h.menuRepo.FindByID(uint(menuItemID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the menu item
	if err := h.menuRepo.Delete(menuItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update positions of remaining items
	if err := h.menuRepo.UpdatePositions(menuItem.Position); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}

// ShowCreateMenuItem displays the menu item creation form
func (h *AdminHandlers) ShowCreateMenuItem(c *gin.Context) {
	pages, err := h.repos.Pages.FindAll()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_create_menu_item.tmpl", h.addCommonData(c, gin.H{
		"title": "Create Menu Item",
		"pages": pages,
	}))
}

// CreateMenuItem handles menu item creation
func (h *AdminHandlers) CreateMenuItem(c *gin.Context) {
	label := c.PostForm("label")
	urlStr := c.PostForm("url")
	pageID := c.PostForm("page_id")

	if label == "" || (urlStr == "" && pageID == "") {
		c.HTML(http.StatusBadRequest, "admin_create_menu_item.tmpl", h.addCommonData(c, gin.H{
			"error": "Label and either URL or Page are required",
			"pages": []models.Page{},
		}))
		return
	}

	menuItem := models.MenuItem{
		Label:    label,
		Position: h.menuRepo.GetNextPosition(),
	}

	// Handle either URL or Page reference
	if pageID != "" {
		pid, _ := utils.ParseUint(pageID)
		menuItem.PageID = &pid
	} else if urlStr != "" {
		menuItem.URL = &urlStr
	}

	if err := h.menuRepo.Create(&menuItem); err != nil {
		pages, err := h.repos.Pages.FindAll()

		if err != nil {
			c.HTML(http.StatusInternalServerError, "admin_create_menu_item.tmpl", h.addCommonData(c, gin.H{
				"error": "Failed to create menu item",
				"item":  menuItem,
				"pages": pages,
			}))
			return
		}
	}

	c.Redirect(http.StatusFound, "/admin/menus")
}

// MoveMenuItem handles menu item reordering
func (h *AdminHandlers) MoveMenuItem(c *gin.Context) {
	id := c.Param("id")
	direction := c.Param("direction")

	menuID, err := utils.ParseUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	if err := h.menuRepo.Move(menuID, direction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item moved successfully"})
}

// ConfirmDeleteMenuItem shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeleteMenuItem(c *gin.Context) {
	id := c.Param("id")

	menuID, err := utils.ParseUint(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	menuItem, err := h.menuRepo.FindByID(menuID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_confirm_delete_menu_item.tmpl", h.addCommonData(c, gin.H{
		"title":    "Confirm Delete Menu Item",
		"menuItem": menuItem,
	}))
}

// EditMenuItem shows the menu item edit form
func (h *AdminHandlers) EditMenuItem(c *gin.Context) {
	var menuItem *models.MenuItem
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	menuID := uint(id)

	if menuItem, err = h.menuRepo.FindByID(menuID); err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	pages, err := h.pageRepo.FindAll()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_edit_menu_item.tmpl", h.addCommonData(c, gin.H{
		"title":    "Edit Menu Item",
		"menuItem": menuItem,
		"pages":    pages,
	}))
}

// UpdateMenuItem handles menu item updates
func (h *AdminHandlers) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")

	menuID, err := utils.ParseUint(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	pages, err := h.pageRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	menuItem, err := h.menuRepo.FindByID(menuID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	menuItem.Label = c.PostForm("label")
	menuItem.URL = nil
	menuItem.PageID = nil

	// Handle either URL or Page reference
	if pageID := c.PostForm("page_id"); pageID != "" {
		pid, _ := utils.ParseUint(pageID)
		menuItem.PageID = &pid
	} else if urlStr := c.PostForm("url"); urlStr != "" {
		menuItem.URL = &urlStr
	}

	if err := h.menuRepo.Update(menuItem); err != nil {

		c.HTML(http.StatusInternalServerError, "admin_edit_menu_item.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update menu item",
			"item":  menuItem,
			"pages": pages,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/menus")
}
