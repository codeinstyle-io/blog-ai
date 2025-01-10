package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/utils"
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
	Timezone    string   `json:"timezone"`
}

type pageRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
	Visible     bool   `json:"visible"`
}

type publishedTime struct {
	Raw            time.Time
	Timezone       string
	UTC            time.Time
	TimezoneOffset int
}

func parseTime(date *string, timezone string) (*publishedTime, error) {
	var parsedTime time.Time

	if date != nil && *date != "" {
		time, err := time.Parse(time.RFC3339, *date)
		if err != nil {
			return nil, err
		}
		parsedTime = time
	} else {
		parsedTime = time.Now()
	}

	// Load the timezone location
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %v", err)
	}

	// Convert the time to the specified timezone
	localTime := parsedTime.In(loc)

	// Calculate timezone offset in minutes
	_, offset := localTime.Zone()
	offsetMinutes := offset / 60 // Convert seconds to minutes

	return &publishedTime{
		Raw:            localTime,
		Timezone:       timezone,
		UTC:            localTime.UTC(),
		TimezoneOffset: offsetMinutes,
	}, nil
}

func (h *AdminHandlers) ApiCreatePost(c *fiber.Ctx) error {
	post := new(postRequest)

	if err := c.BodyParser(post); err != nil {
		// TODO: Log error
		fmt.Printf("Failed to parse request body: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get the logged in user
	exists := c.Locals("user")
	if exists == nil {
		flash.Error(c, "User session not found")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "User session not found",
		})
	}
	user := exists.(*models.User)

	publishedAt, err := parseTime(post.PublishedAt, post.Timezone)
	if err != nil {
		// TODO: Log error
		fmt.Printf("Failed to parse publishedAt: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid publish date"})
	}

	newPost := &models.Post{
		Title:                     post.Title,
		Slug:                      post.Slug,
		Content:                   post.Content,
		Excerpt:                   &post.Excerpt,
		Visible:                   post.Visible,
		PublishedAt:               publishedAt.Raw,
		PublishedAtTimezone:       publishedAt.Timezone,
		PublishedAtUTC:            publishedAt.UTC,
		PublishedAtTimeZoneOffset: publishedAt.TimezoneOffset,
		AuthorID:                  user.ID,
	}

	if err := h.repos.Posts.Create(newPost); err != nil {
		// TODO: Log error
		fmt.Printf("Error creating post: %v\n", err)

		if utils.IsConstraintError(err) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Post with the same slug already exists"})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create post"})
	}

	if err := h.repos.Posts.AssociateTags(newPost, post.Tags); err != nil {
		// TODO: Log error
		fmt.Printf("Error associating tags: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to associate tags"})
	}

	flash.Success(c, "Post created successfully")

	return c.JSON(fiber.Map{"message": "Post created successfully", "redirect": "/admin/posts"})
}

func (h *AdminHandlers) ApiUpdatePost(c *fiber.Ctx) error {
	post := new(postRequest)

	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	postToUpdate, err := h.repos.Posts.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Post not found"})
	}

	if err := c.BodyParser(post); err != nil {
		// TODO: Log error
		fmt.Printf("Failed to parse request body: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	publishedAt, err := parseTime(post.PublishedAt, post.Timezone)
	if err != nil {
		// TODO: Log error
		fmt.Printf("Failed to parse publishedAt: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid publish date"})
	}

	postToUpdate.Title = post.Title
	postToUpdate.Slug = post.Slug
	postToUpdate.Content = post.Content
	postToUpdate.Excerpt = &post.Excerpt
	postToUpdate.Visible = post.Visible
	postToUpdate.PublishedAt = publishedAt.Raw
	postToUpdate.PublishedAtTimezone = publishedAt.Timezone
	postToUpdate.PublishedAtUTC = publishedAt.UTC
	postToUpdate.PublishedAtTimeZoneOffset = publishedAt.TimezoneOffset

	if err := h.repos.Posts.Update(postToUpdate); err != nil {
		// TODO: Log error
		fmt.Printf("Error updating post: %v\n", err)

		if utils.IsConstraintError(err) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Post with the same slug already exists"})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update post"})
	}

	if err := h.repos.Posts.AssociateTags(postToUpdate, post.Tags); err != nil {
		// TODO: Log error
		fmt.Printf("Error associating tags: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to associate tags"})
	}

	flash.Success(c, "Post updated successfully")

	return c.JSON(fiber.Map{"message": "Post updated successfully", "redirect": "/admin/posts"})
}

func (h *AdminHandlers) ApiCreatePage(c *fiber.Ctx) error {
	page := new(pageRequest)
	if err := c.BodyParser(page); err != nil {
		// TODO: Log error
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	newPage := &models.Page{
		Title:       page.Title,
		Slug:        page.Slug,
		Content:     page.Content,
		ContentType: page.ContentType,
		Visible:     page.Visible,
	}

	if err := h.repos.Pages.Create(newPage); err != nil {
		// TODO: Log error
		fmt.Printf("Error creating page: %v\n", err)

		if utils.IsConstraintError(err) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Page with the same slug already exists"})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create page"})
	}

	flash.Success(c, "Page created successfully")

	return c.JSON(fiber.Map{"message": "Page created successfully", "redirect": "/admin/pages"})
}

func (h *AdminHandlers) ApiUpdatePage(c *fiber.Ctx) error {
	page := new(pageRequest)
	if err := c.BodyParser(page); err != nil {
		// TODO: Log error
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page ID"})
	}

	pageToUpdate, err := h.repos.Pages.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Page not found"})
	}

	pageToUpdate.Title = page.Title
	pageToUpdate.Slug = page.Slug
	pageToUpdate.Content = page.Content
	pageToUpdate.ContentType = page.ContentType
	pageToUpdate.Visible = page.Visible

	if err := h.repos.Pages.Update(pageToUpdate); err != nil {
		// TODO: Log error
		fmt.Printf("Error updating page: %v\n", err)

		if utils.IsConstraintError(err) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Page with the same slug already exists"})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update page"})
	}

	flash.Success(c, "Page updated successfully")

	return c.JSON(fiber.Map{"message": "Page updated successfully", "redirect": "/admin/pages"})
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
