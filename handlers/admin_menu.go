package handlers

import (
	"net/http"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
)

// ListMenuItems shows all menu items
func (h *AdminHandlers) ListMenuItems(c *gin.Context) {
	var menuItems []db.MenuItem
	if err := h.db.Preload("Page").Order("position").Find(&menuItems).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	lastPosition := len(menuItems) - 1
	c.HTML(http.StatusOK, "admin_menu_items.tmpl", gin.H{
		"title":        "Menu Items",
		"menuItems":    menuItems,
		"lastPosition": lastPosition,
	})
}

// ShowCreateMenuItem displays the menu item creation form
func (h *AdminHandlers) ShowCreateMenuItem(c *gin.Context) {
	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_create_menu_item.tmpl", gin.H{
		"title": "Create Menu Item",
		"pages": pages,
	})
}

// CreateMenuItem handles menu item creation
func (h *AdminHandlers) CreateMenuItem(c *gin.Context) {
	label := c.PostForm("label")
	urlStr := c.PostForm("url")
	pageID := c.PostForm("page_id")

	if label == "" || (urlStr == "" && pageID == "") {
		c.HTML(http.StatusBadRequest, "admin_create_menu_item.tmpl", gin.H{
			"error": "Label and either URL or Page are required",
			"pages": []db.Page{},
		})
		return
	}

	menuItem := db.MenuItem{
		Label:    label,
		Position: h.getNextMenuPosition(),
	}

	// Handle either URL or Page reference
	if pageID != "" {
		pid := parseUint(pageID)
		menuItem.PageID = &pid
	} else if urlStr != "" {
		menuItem.URL = &urlStr
	}

	if err := h.db.Create(&menuItem).Error; err != nil {
		var pages []db.Page
		h.db.Find(&pages)
		c.HTML(http.StatusInternalServerError, "admin_create_menu_item.tmpl", gin.H{
			"error": "Failed to create menu item",
			"item":  menuItem,
			"pages": pages,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/menus")
}

// MoveMenuItem handles menu item reordering
func (h *AdminHandlers) MoveMenuItem(c *gin.Context) {
	id := c.Param("id")
	direction := c.Param("direction")

	// Start transaction
	tx := h.db.Begin()

	var currentItem db.MenuItem
	if err := tx.First(&currentItem, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		return
	}

	// Find adjacent item
	var adjacentItem db.MenuItem
	if direction == "up" {
		if err := tx.Where("position < ?", currentItem.Position).Order("position DESC").First(&adjacentItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item already at top"})
			return
		}
	} else {
		if err := tx.Where("position > ?", currentItem.Position).Order("position ASC").First(&adjacentItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item already at bottom"})
			return
		}
	}

	// Swap positions
	currentPos := currentItem.Position
	adjacentPos := adjacentItem.Position

	if err := tx.Model(&currentItem).Update("position", adjacentPos).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update position"})
		return
	}

	if err := tx.Model(&adjacentItem).Update("position", currentPos).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update position"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Menu item moved successfully"})
}

// ConfirmDeleteMenuItem shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	var menuItem db.MenuItem
	if err := h.db.First(&menuItem, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_confirm_delete_menu_item.tmpl", gin.H{
		"title":    "Confirm Delete Menu Item",
		"menuItem": menuItem,
	})
}

// DeleteMenuItem removes a menu item
func (h *AdminHandlers) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	var menuItem db.MenuItem
	if err := h.db.First(&menuItem, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	// Start transaction
	tx := h.db.Begin()

	// Delete the menu item
	if err := tx.Delete(&menuItem).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	// Update positions of remaining items
	if err := tx.Model(&db.MenuItem{}).Where("position > ?", menuItem.Position).
		UpdateColumn("position", h.db.Raw("position - 1")).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}

// EditMenuItem shows the menu item edit form
func (h *AdminHandlers) EditMenuItem(c *gin.Context) {
	id := c.Param("id")
	var menuItem db.MenuItem
	if err := h.db.First(&menuItem, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_edit_menu_item.tmpl", gin.H{
		"title":    "Edit Menu Item",
		"menuItem": menuItem,
		"pages":    pages,
	})
}

// UpdateMenuItem handles menu item updates
func (h *AdminHandlers) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	var menuItem db.MenuItem
	if err := h.db.First(&menuItem, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	menuItem.Label = c.PostForm("title")

	// Reset both URL and PageID
	menuItem.URL = nil
	menuItem.PageID = nil

	// Handle either URL or Page reference
	if pageID := c.PostForm("page_id"); pageID != "" {
		pid := parseUint(pageID)
		menuItem.PageID = &pid
	} else if urlStr := c.PostForm("url"); urlStr != "" {
		menuItem.URL = &urlStr
	}

	if err := h.db.Save(&menuItem).Error; err != nil {
		var pages []db.Page
		h.db.Find(&pages)
		c.HTML(http.StatusInternalServerError, "admin_edit_menu_item.tmpl", gin.H{
			"error": "Failed to update menu item",
			"item":  menuItem,
			"pages": pages,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/menu")
}
