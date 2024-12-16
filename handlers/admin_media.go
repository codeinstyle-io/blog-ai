package handlers

import (
	"fmt"
	"net/http"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
)

// AdminMediaHandlers handles media routes
type AdminMediaHandlers struct {
	*BaseHandler
	storage   storage.Provider
	mediaRepo models.MediaRepository
}

// NewAdminMediaHandlers creates a new AdminMediaHandlers instance
func NewAdminMediaHandlers(repos *repository.Repositories, config *config.Config, storage storage.Provider) *AdminMediaHandlers {
	return &AdminMediaHandlers{
		BaseHandler: NewBaseHandler(repos, config),
		storage:     storage,
		mediaRepo:   repos.Media,
	}
}

// ListMedia displays the list of media files
func (h *AdminMediaHandlers) ListMedia(c *gin.Context) {
	media, err := h.mediaRepo.FindAll()
	if err != nil {
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
func (h *AdminMediaHandlers) ShowUploadMedia(c *gin.Context) {
	data := gin.H{
		"title": "Upload Media",
	}
	data = h.addCommonData(c, data)
	c.HTML(http.StatusOK, "admin_media_upload.tmpl", data)
}

// UploadMedia handles media file upload
func (h *AdminMediaHandlers) UploadMedia(c *gin.Context) {
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
	media := &models.Media{
		Name:        file.Filename,
		Path:        filename,
		Size:        file.Size,
		Description: description,
	}

	err = h.mediaRepo.Create(media)
	if err != nil {
		// Clean up file if database insert fails
		if err := h.storage.Delete(filename); err != nil {
			c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", h.addCommonData(c, gin.H{
				"error": fmt.Sprintf("Failed to delete file: %v", err),
			}))
			return
		}

		c.HTML(http.StatusInternalServerError, "admin_media_upload.tmpl", h.addCommonData(c, gin.H{
			"error": fmt.Sprintf("Failed to save media record: %v", err),
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/media")
}

// DeleteMedia handles media deletion
func (h *AdminMediaHandlers) DeleteMedia(c *gin.Context) {
	mediaID, err := utils.ParseUint(c.Param("id"))

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	media, err := h.mediaRepo.FindByID(mediaID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Delete file using storage provider
	if err := h.storage.Delete(media.Path); err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Delete record
	if err := h.mediaRepo.Delete(media); err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

// GetMediaList returns a JSON list of media for AJAX requests
func (h *AdminMediaHandlers) GetMediaList(c *gin.Context) {
	media, err := h.mediaRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch media"})
		return
	}

	c.JSON(http.StatusOK, media)
}

// ConfirmDeleteMedia shows the delete confirmation page
func (h *AdminMediaHandlers) ConfirmDeleteMedia(c *gin.Context) {
	mediaID, err := utils.ParseUint(c.Param("id"))

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	media, err := h.mediaRepo.FindByID(mediaID)
	if err != nil {
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

func (h *AdminMediaHandlers) addCommonData(c *gin.Context, data gin.H) gin.H {
	settings, _ := h.repos.Settings.Get()

	data["settings"] = settings
	data["user"] = c.MustGet("user").(*models.User)
	return data
}
