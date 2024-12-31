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
}

type pageRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
	Visible     bool   `json:"visible"`
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
	settings := c.Locals("settings").(*models.Settings)

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
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	publishedAt, err := parseTime(post.PublishedAt, settings.Timezone)
	if err != nil {
		// TODO: Log error
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid publishedAt"})
	}

	postToUpdate.Title = post.Title
	postToUpdate.Slug = post.Slug
	postToUpdate.Content = post.Content
	postToUpdate.Excerpt = &post.Excerpt
	postToUpdate.Visible = post.Visible
	postToUpdate.PublishedAt = *publishedAt

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
