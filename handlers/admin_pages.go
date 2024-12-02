package handlers

import (
	"net/http"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
)

// ListPages shows all pages
func (h *AdminHandlers) ListPages(c *gin.Context) {
	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_pages.tmpl", gin.H{
		"title": "Pages",
		"pages": pages,
	})
}

// ShowCreatePage displays the page creation form
func (h *AdminHandlers) ShowCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_page.tmpl", gin.H{
		"title": "Create Page",
	})
}

// CreatePage handles page creation
func (h *AdminHandlers) CreatePage(c *gin.Context) {
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	visible := c.PostForm("visible") == "on"

	if title == "" || slug == "" || content == "" {
		c.HTML(http.StatusBadRequest, "admin_create_page.tmpl", gin.H{
			"error": "All fields are required",
		})
		return
	}

	page := db.Page{
		Title:   title,
		Slug:    slug,
		Content: content,
		Visible: visible,
	}

	if err := h.db.Create(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_page.tmpl", gin.H{
			"error": "Failed to create page",
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/pages")
}

// EditPage shows the page edit form
func (h *AdminHandlers) EditPage(c *gin.Context) {
	id := c.Param("id")
	var page db.Page
	if err := h.db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_edit_page.tmpl", gin.H{
		"title": "Edit Page",
		"page":  page,
	})
}

// UpdatePage handles page updates
func (h *AdminHandlers) UpdatePage(c *gin.Context) {
	id := c.Param("id")
	var page db.Page
	if err := h.db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	visible := c.PostForm("visible") == "on"

	if title == "" || slug == "" || content == "" {
		c.HTML(http.StatusBadRequest, "admin_edit_page.tmpl", gin.H{
			"error": "All fields are required",
			"page":  page,
		})
		return
	}

	page.Title = title
	page.Slug = slug
	page.Content = content
	page.Visible = visible

	if err := h.db.Save(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_page.tmpl", gin.H{
			"error": "Failed to update page",
			"page":  page,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/pages")
}

// ConfirmDeletePage shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePage(c *gin.Context) {
	id := c.Param("id")
	var page db.Page
	if err := h.db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_confirm_delete_page.tmpl", gin.H{
		"title": "Confirm Delete Page",
		"page":  page,
	})
}

// DeletePage removes a page
func (h *AdminHandlers) DeletePage(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&db.Page{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
