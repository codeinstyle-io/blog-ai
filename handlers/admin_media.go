package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
)

// ListMedia displays the list of media files
func (h *AdminHandlers) ListMedia(c *gin.Context) {
	var media []db.Media
	result := h.db.Order("created_at desc").Find(&media)
	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	data := gin.H{
		"title": "Media Library",
		"media": media,
	}
	data = h.addCommonData(c, data)
	c.HTML(http.StatusOK, "admin_media_list.tmpl", data)
}

// ShowUploadMedia displays the upload media form
func (h *AdminHandlers) ShowUploadMedia(c *gin.Context) {
	data := gin.H{
		"title": "Upload Media",
	}
	data = h.addCommonData(c, data)
	c.HTML(http.StatusOK, "admin_media_upload.tmpl", data)
}

// UploadMedia handles media file upload
func (h *AdminHandlers) UploadMedia(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_media_upload.tmpl", h.addCommonData(c, gin.H{
			"error": "No file uploaded",
		}))
		return
	}

	description := c.PostForm("description")

	// Save file using storage provider
	filename, err := h.storage.Save(file)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", h.addCommonData(c, gin.H{
			"error": fmt.Sprintf("Failed to save file: %v", err),
		}))
		return
	}

	// Create media record
	media := db.Media{
		Name:        file.Filename,
		Path:        filename,
		Size:        file.Size,
		Description: description,
	}

	if result := h.db.Create(&media); result.Error != nil {
		// Clean up file if database insert fails
		if err := h.storage.Delete(filename); err != nil {
			c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", h.addCommonData(c, gin.H{
				"error": fmt.Sprintf("Failed to delete file: %v", err),
			}))
			return
		}

		c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", h.addCommonData(c, gin.H{
			"error": fmt.Sprintf("Failed to save media record: %v", result.Error),
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/media")
}

// DeleteMedia handles media deletion
func (h *AdminHandlers) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	mediaID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	var media db.Media
	if err := h.db.First(&media, mediaID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Delete file using storage provider
	if err := h.storage.Delete(media.Path); err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Delete record
	if err := h.db.Delete(&media).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

// GetMediaList returns a JSON list of media for AJAX requests
func (h *AdminHandlers) GetMediaList(c *gin.Context) {
	var media []db.Media
	if err := h.db.Order("created_at desc").Find(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch media"})
		return
	}

	c.JSON(http.StatusOK, media)
}

// ConfirmDeleteMedia shows the delete confirmation page
func (h *AdminHandlers) ConfirmDeleteMedia(c *gin.Context) {
	id := c.Param("id")
	mediaID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	var media db.Media
	if err := h.db.First(&media, mediaID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	data := gin.H{
		"title": "Delete Media",
		"media": media,
	}
	data = h.addCommonData(c, data)
	c.HTML(http.StatusOK, "admin_confirm_delete_media.tmpl", data)
}
