package handlers

import (
	"net/http"

	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
)

// ListPages shows all pages
func (h *AdminHandlers) ListPages(c *gin.Context) {
	pages, err := h.pageRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}
	c.HTML(http.StatusOK, "admin_pages.tmpl", h.addCommonData(c, gin.H{
		"title": "Pages",
		"pages": pages,
	}))
}

// ShowCreatePage displays the page creation form
func (h *AdminHandlers) ShowCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_page.tmpl", h.addCommonData(c, gin.H{
		"title": "Create Page",
	}))
}

// CreatePage handles page creation
func (h *AdminHandlers) CreatePage(c *gin.Context) {
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	visible := c.PostForm("visible") == "on"

	if title == "" || slug == "" || content == "" {
		c.HTML(http.StatusBadRequest, "admin_create_page.tmpl", h.addCommonData(c, gin.H{
			"error": "All fields are required",
		}))
		return
	}

	page := models.Page{
		Title:   title,
		Slug:    slug,
		Content: content,
		Visible: visible,
	}

	if err := h.pageRepo.Create(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_page.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to create page",
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/pages")
}

// EditPage shows the page edit form
func (h *AdminHandlers) EditPage(c *gin.Context) {
	id := c.Param("id")
	pageID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	page, err := h.pageRepo.FindByID(pageID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_edit_page.tmpl", h.addCommonData(c, gin.H{
		"title": "Edit Page",
		"page":  page,
	}))
}

// UpdatePage handles page updates
func (h *AdminHandlers) UpdatePage(c *gin.Context) {
	id := c.Param("id")
	pageID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	page, err := h.pageRepo.FindByID(pageID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	visible := c.PostForm("visible") == "on"

	if title == "" || slug == "" || content == "" {
		c.HTML(http.StatusBadRequest, "admin_edit_page.tmpl", h.addCommonData(c, gin.H{
			"error": "All fields are required",
			"page":  page,
		}))
		return
	}

	page.Title = title
	page.Slug = slug
	page.Content = content
	page.Visible = visible

	if err := h.pageRepo.Update(page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_page.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update page",
			"page":  page,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/pages")
}

// ConfirmDeletePage shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePage(c *gin.Context) {
	id := c.Param("id")
	pageID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	page, err := h.pageRepo.FindByID(pageID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_confirm_delete_page.tmpl", h.addCommonData(c, gin.H{
		"title": "Confirm Delete Page",
		"page":  page,
	}))
}

// DeletePage removes a page
func (h *AdminHandlers) DeletePage(c *gin.Context) {
	var menuItemCount int64
	id := c.Param("id")
	pageId, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	page, err := h.pageRepo.FindByID(pageId)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Check for menu item references

	err = h.pageRepo.CountRelatedMenuItems(page.ID, &menuItemCount)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	if menuItemCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete page: it is referenced by one or more menu items. Please remove the menu items first.",
		})
		return
	}

	// If no references exist, proceed with deletion
	if err := h.pageRepo.Delete(page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
