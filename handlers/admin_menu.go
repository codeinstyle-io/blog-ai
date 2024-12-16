package handlers

import (
	"net/http"
	"strconv"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/models"
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
	var menuItems []struct {
		Title    string `json:"title"`
		URL      string `json:"url"`
		Position int    `json:"position"`
	}

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
	for _, item := range menuItems {
		menuItem := &db.MenuItem{
			Label:    item.Title,
			URL:      &item.URL,
			Position: item.Position,
		}

		if err := h.menuRepo.Create(menuItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu items saved successfully"})
}

// DeleteMenuItem deletes a menu item
func (h *AdminHandlers) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	menuItemID, err := strconv.ParseUint(id, 10, 32)
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
	if err := h.menuRepo.Delete(uint(menuItemID)); err != nil {
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
	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
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
			"pages": []db.Page{},
		}))
		return
	}

	menuItem := models.MenuItem{
		Label:    label,
		Position: h.menuRepo.GetNextPosition(),
	}

	// Handle either URL or Page reference
	if pageID != "" {
		pid := parseUint(pageID)
		menuItem.PageID = &pid
	} else if urlStr != "" {
		menuItem.URL = &urlStr
	}

	if err := h.menuRepo.Create(&menuItem); err != nil {
		var pages []db.Page
		h.db.Find(&pages)
		c.HTML(http.StatusInternalServerError, "admin_create_menu_item.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to create menu item",
			"item":  menuItem,
			"pages": pages,
		}))
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
	if err := h.menuRepo.Find(id, &menuItem); err != nil {
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

	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	menuID := uint(id)

	menuItem, err := h.menuRepo.FindByID(menuID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	menuItem.Label = c.PostForm("label")

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

	if err := h.menuRepo.Save(&menuItem); err != nil {
		var pages []db.Page
		h.db.Find(&pages)
		c.HTML(http.StatusInternalServerError, "admin_edit_menu_item.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update menu item",
			"item":  menuItem,
			"pages": pages,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/menus")
}
