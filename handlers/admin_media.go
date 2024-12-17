package handlers

import (
	"fmt"
	"net/http"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"codeinstyle.io/captain/utils"
	"github.com/gofiber/fiber/v2"
)

// AdminMediaHandlers handles media routes
type AdminMediaHandlers struct {
	*BaseHandlers
	storage   storage.Provider
	mediaRepo models.MediaRepository
}

// NewAdminMediaHandlers creates a new AdminMediaHandlers instance
func NewAdminMediaHandlers(repos *repository.Repositories, config *config.Config, storage storage.Provider) *AdminMediaHandlers {
	return &AdminMediaHandlers{
		BaseHandlers: NewBaseHandlers(repos, config),
		storage:      storage,
		mediaRepo:    repos.Media,
	}
}

// ListMedia displays the list of media files
func (h *AdminMediaHandlers) ListMedia(c *fiber.Ctx) error {
	media, err := h.mediaRepo.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", h.addCommonData(c, fiber.Map{}))
	}

	return c.Render("admin_media_list", h.addCommonData(c, fiber.Map{
		"title": "Media Library",
		"media": media,
	}))
}

// ShowUploadMedia displays the upload media form
func (h *AdminMediaHandlers) ShowUploadMedia(c *fiber.Ctx) error {
	return c.Render("admin_media_upload", h.addCommonData(c, fiber.Map{
		"title": "Upload Media",
	}))
}

// UploadMedia handles media file upload
func (h *AdminMediaHandlers) UploadMedia(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_media_upload", h.addCommonData(c, fiber.Map{
			"error": "No file uploaded",
		}))
	}

	description := c.FormValue("description")

	// Save file using storage provider
	filename, err := h.storage.Save(file)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_media_upload", h.addCommonData(c, fiber.Map{
			"error": fmt.Sprintf("Failed to save file: %v", err),
		}))
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
			return c.Status(http.StatusInternalServerError).Render("admin_media_upload", h.addCommonData(c, fiber.Map{
				"error": fmt.Sprintf("Failed to delete file: %v", err),
			}))
		}

		return c.Status(http.StatusInternalServerError).Render("admin_media_upload", h.addCommonData(c, fiber.Map{
			"error": fmt.Sprintf("Failed to save media record: %v", err),
		}))
	}

	return c.Redirect("/admin/media")
}

// DeleteMedia handles media deletion
func (h *AdminMediaHandlers) DeleteMedia(c *fiber.Ctx) error {
	mediaID, err := utils.ParseUint(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", h.addCommonData(c, fiber.Map{}))
	}

	media, err := h.mediaRepo.FindByID(mediaID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", h.addCommonData(c, fiber.Map{}))
	}

	// Delete file using storage provider
	if err := h.storage.Delete(media.Path); err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", h.addCommonData(c, fiber.Map{}))
	}

	// Delete record
	if err := h.mediaRepo.Delete(media); err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", h.addCommonData(c, fiber.Map{}))
	}

	return c.JSON(fiber.Map{"message": "Media deleted successfully"})
}

// GetMediaList returns a JSON list of media for AJAX requests
func (h *AdminMediaHandlers) GetMediaList(c *fiber.Ctx) error {
	media, err := h.mediaRepo.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch media"})
	}

	return c.JSON(media)
}

// ConfirmDeleteMedia shows the delete confirmation page
func (h *AdminMediaHandlers) ConfirmDeleteMedia(c *fiber.Ctx) error {
	mediaID, err := utils.ParseUint(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", h.addCommonData(c, fiber.Map{}))
	}

	media, err := h.mediaRepo.FindByID(mediaID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", h.addCommonData(c, fiber.Map{}))
	}

	return c.Render("admin_confirm_delete_media", h.addCommonData(c, fiber.Map{
		"title": "Delete Media",
		"media": media,
	}))
}

func (h *AdminMediaHandlers) addCommonData(c *fiber.Ctx, data fiber.Map) fiber.Map {
	settings, _ := h.repos.Settings.Get()

	data["settings"] = settings
	data["user"] = c.Locals("user").(*models.User)
	return data
}
