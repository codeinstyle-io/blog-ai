package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
)

// ListMedia displays the list of media files
func (h *AdminHandlers) ListMedia(c *gin.Context) {
	var media []db.Media
	result := h.db.Order("created_at desc").Find(&media)
	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
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
		c.HTML(http.StatusBadRequest, "admin_media_upload.tmpl", gin.H{
			"error": "No file uploaded",
		})
		return
	}

	description := c.PostForm("description")

	// Create media directory if it doesn't exist
	mediaDir := "./media"
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", gin.H{
			"error": fmt.Sprintf("Failed to create media directory: %v", err),
		})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().Unix(), strings.TrimSuffix(file.Filename, ext), ext)
	filepath := filepath.Join(mediaDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", gin.H{
			"error": fmt.Sprintf("Failed to save file: %v", err),
		})
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
		os.Remove(filepath)
		c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", gin.H{
			"error": fmt.Sprintf("Failed to save media record: %v", result.Error),
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/media")
}

// DeleteMedia handles media deletion
func (h *AdminHandlers) DeleteMedia(c *gin.Context) {
	id := c.Param("id")
	mediaID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", gin.H{})
		return
	}

	var media db.Media
	if err := h.db.First(&media, mediaID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	// Delete file
	filepath := filepath.Join("./media", media.Path)
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	// Delete record
	if err := h.db.Delete(&media).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.Redirect(http.StatusFound, "/admin/media")
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
		c.HTML(http.StatusBadRequest, "500.tmpl", gin.H{})
		return
	}

	var media db.Media
	if err := h.db.First(&media, mediaID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	data := gin.H{
		"title": "Delete Media",
		"media": media,
	}
	data = h.addCommonData(c, data)
	c.HTML(http.StatusOK, "admin_media_delete.tmpl", data)
}
