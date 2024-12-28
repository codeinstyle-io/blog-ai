package handlers

import (
	"net/http"

	"github.com/captain-corp/captain/models"
	"github.com/gofiber/fiber/v2"
)

// tagResponse struct for API responses
type tagResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func (h *AdminHandlers) ApiCreatePost(c *fiber.Ctx) error {
	post := new(models.Post)

	if err := c.BodyParser(post); err != nil {
		// TODO: Log error
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.repos.Posts.Create(post); err != nil {
		// TODO: Log error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.JSON(fiber.Map{"message": "Post created successfully"})
}

func (h *AdminHandlers) ApiUpdatePost(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
}

func (h *AdminHandlers) ApiGetPage(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
}

func (h *AdminHandlers) ApiCreatePage(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
}

func (h *AdminHandlers) ApiUpdatePage(c *fiber.Ctx) error {
	return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
}

// ApiGetMediaList returns a JSON list of media for AJAX requests
func (h *AdminMediaHandlers) ApiGetMediaList(c *fiber.Ctx) error {
	media, err := h.mediaRepo.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch media"})
	}

	return c.JSON(media)
}

// ApiGetTags returns a list of tags for API consumption
func (h *AdminHandlers) ApiGetTags(c *fiber.Ctx) error {
	tags, err := h.repos.Tags.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load tags",
		})
	}

	var response []tagResponse
	for _, tag := range tags {
		response = append(response, tagResponse{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}

	return c.JSON(response)
}
