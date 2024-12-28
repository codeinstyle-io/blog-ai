package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/captain-corp/captain/models"
	"github.com/gofiber/fiber/v2"
)

// tagResponse struct for API responses
type tagResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type postRequest struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Content     string   `json:"content"`
	Excerpt     string   `json:"excerpt"`
	Tags        []string `json:"tags"`
	Visible     bool     `json:"visible"`
	PublishedAt *string  `json:"publishedAt"`
}

func parseTime(date *string, timezone string) (*time.Time, error) {
	var parsedTime time.Time

	loc, err := time.LoadLocation(timezone)

	if err != nil {
		return nil, err
	}

	if date != nil {
		time, err := time.Parse(time.RFC3339, *date)
		if err != nil {
			// TODO: Log error
			return nil, err
		}
		parsedTime = time
	} else {
		parsedTime = time.Now()
	}

	parsedTime = parsedTime.In(loc)

	return &parsedTime, nil
}

func (h *AdminHandlers) ApiCreatePost(c *fiber.Ctx) error {
	post := new(postRequest)
	settings := c.Locals("settings").(*models.Settings)

	if err := c.BodyParser(post); err != nil {
		// TODO: Log error
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	publishedAt, err := parseTime(post.PublishedAt, settings.Timezone)
	if err != nil {
		// TODO: Log error
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid publishedAt"})
	}

	newPost := &models.Post{
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		Excerpt:     &post.Excerpt,
		Visible:     post.Visible,
		PublishedAt: *publishedAt,
		AuthorID:    1, //TODO: Get logged in user
	}

	if err := h.repos.Posts.Create(newPost); err != nil {
		// TODO: Log error
		fmt.Printf("Error creating post: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
	}

	if err := h.repos.Posts.AssociateTags(newPost, post.Tags); err != nil {
		// TODO: Log error
		fmt.Printf("Error associating tags: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to associate tags"})
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
