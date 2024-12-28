package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/utils"

	"github.com/gofiber/fiber/v2"
)

// ListPosts shows all posts for admin
func (h *AdminHandlers) ListPosts(c *fiber.Ctx) error {
	posts, err := h.repos.Posts.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert post times to user timezone
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(loc)
	}

	return c.Render("admin_posts", fiber.Map{
		"title": "Posts",
		"posts": posts,
	})
}

// ShowCreatePost displays the post creation form
func (h *AdminHandlers) ShowCreatePost(c *fiber.Ctx) error {
	return c.Render("admin_create_post", fiber.Map{
		"title": "Create Post",
	})
}

// CreatePost handles post creation
func (h *AdminHandlers) CreatePost(c *fiber.Ctx) error {

	// Get the logged in user
	exists := c.Locals("user")
	if exists == nil {
		flash.Error(c, "User session not found")
		return c.Status(http.StatusInternalServerError).Render("admin_create_post", fiber.Map{
			"title": "Posts",
		})
	}
	user := exists.(*models.User)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		flash.Error(c, "Invalid form data")
		return c.Status(http.StatusBadRequest).Render("admin_create_post", fiber.Map{
			"title": "Posts",
			"post":  &models.Post{},
		})
	}

	// Parse form data
	post, err := postFromForm(form)
	post.AuthorID = user.ID

	if err != nil {
		fmt.Printf("Error parsing form: %v\n", err)
		flash.Error(c, "Unable to parse form into post")
		return c.Status(http.StatusBadRequest).Render("admin_create_post", fiber.Map{
			"title": "Posts",
			"post":  &post,
		})
	}

	// Basic validation
	if post.Title == "" || post.Slug == "" || post.Content == "" {
		flash.Error(c, "Title, slug and content are required")
		return c.Status(http.StatusBadRequest).Render("admin_create_post", fiber.Map{
			"title": "Posts",
			"post":  &post,
		})
	}

	// Handle tags
	tags := strings.Split(c.FormValue("tags"), ",")

	// Create post with transaction to ensure atomic operation
	if err := h.repos.Posts.Create(post); err != nil {
		flash.Error(c, "Failed to create post")
		return c.Status(http.StatusInternalServerError).Render("admin_create_post", fiber.Map{
			"post": &post,
		})
	}

	if err := h.repos.Posts.AssociateTags(post, tags); err != nil {
		fmt.Printf("Error associating tags: %v\n", err)
		flash.Error(c, "Failed to associate tags")
		return c.Status(http.StatusInternalServerError).Render("admin_create_post", fiber.Map{
			"title": "Posts",
			"post":  &post,
		})
	}

	flash.Success(c, "Post created successfully")
	return c.Redirect("/admin/posts")
}

func (h *AdminHandlers) UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	postID, err := utils.ParseUint(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	post, err := h.repos.Posts.FindByID(postID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	// Parse form data
	post.Title = c.FormValue("title")
	post.Slug = c.FormValue("slug")
	post.Content = c.FormValue("content")
	post.Visible = c.FormValue("visible") == "on"
	excerpt := c.FormValue("excerpt")

	if excerpt != "" {
		post.Excerpt = &excerpt
	}

	if err != nil {
		flash.Error(c, "Unable to parse form into post")
		return c.Status(http.StatusBadRequest).Render("admin_edit_post", fiber.Map{
			"title": "Posts",
			"post":  &post,
		})
	}

	// Basic validation
	if post.Title == "" || post.Slug == "" || post.Content == "" {
		flash.Error(c, "Title, slug and content are required")
		return c.Status(http.StatusBadRequest).Render("admin_edit_post", fiber.Map{
			"title": "Posts",
			"post":  post,
		})
	}

	// Handle tags
	tags := strings.Split(c.FormValue("tags"), ",")

	// Update post with transaction to ensure atomic operation
	if err := h.repos.Posts.Update(post); err != nil {
		flash.Error(c, "Unable to update post")
		return c.Status(http.StatusInternalServerError).Render("admin_edit_post", fiber.Map{
			"title": "Posts",
			"post":  post,
		})
	}

	if err := h.repos.Posts.AssociateTags(post, tags); err != nil {
		flash.Error(c, "Failed to associate tags")
		return c.Status(http.StatusInternalServerError).Render("admin_edit_post", fiber.Map{
			"title": "Posts",
			"post":  post,
		})
	}

	return c.Redirect("/admin/posts")
}

// ListPostsByTag shows all posts for a specific tag
func (h *AdminHandlers) ListPostsByTag(c *fiber.Ctx) error {
	id := c.Params("id")

	tagID, err := utils.ParseUint(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Find the tag
	tag, err := h.repos.Tags.FindByID(tagID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	posts, err := h.repos.Posts.FindByTag(tag.Slug)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert post times to user timezone
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(loc)
	}

	return c.Render("admin_tag_posts", fiber.Map{
		"title": fmt.Sprintf("Posts tagged with '%s'", tag.Name),
		"posts": posts,
		"tag":   tag,
	})
}

// ConfirmDeletePost shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	tagID, err := utils.ParseUint(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	post, err := h.repos.Posts.FindByID(tagID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_post", fiber.Map{
		"title": "Confirm Post deletion",
		"post":  post,
	})
}

// DeletePost handles post deletion
func (h *AdminHandlers) DeletePost(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid post ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid post ID",
			"redirect": "/admin/posts",
		})
	}

	post, err := h.repos.Posts.FindByID(id)
	if err != nil {
		flash.Error(c, "Post not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":    "Post not found",
			"redirect": "/admin/posts",
		})
	}

	// Delete post
	if err := h.repos.Posts.Delete(post); err != nil {
		flash.Error(c, "Failed to delete post")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to delete post",
			"redirect": "/admin/posts",
		})
	}

	flash.Success(c, "Post deleted successfully")
	return c.JSON(fiber.Map{
		"message":  "Post deleted successfully",
		"redirect": "/admin/posts",
	})
}

func (h *AdminHandlers) EditPost(c *fiber.Ctx) error {
	var allTags []*models.Tag
	id := c.Params("id")
	tagID, err := utils.ParseUint(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	post, err := h.repos.Posts.FindByID(tagID)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	allTags, err = h.repos.Tags.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert post time to user timezone
	post.PublishedAt = post.PublishedAt.In(loc)

	return c.Render("admin_edit_post", fiber.Map{
		"title":    "Edit Post",
		"post":     post,
		"allTags":  allTags,
		"postTags": post.Tags,
	})
}

func postFromForm(form *multipart.Form) (*models.Post, error) {
	post := &models.Post{
		Title:   form.Value["title"][0],
		Slug:    form.Value["slug"][0],
		Content: form.Value["content"][0],
		Excerpt: &form.Value["excerpt"][0],
	}

	var err error

	return post, err
}
