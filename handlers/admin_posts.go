package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
)

// ListPosts shows all posts for admin
func (h *AdminHandlers) ListPosts(c *gin.Context) {
	posts, err := h.repos.Posts.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
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

	c.HTML(http.StatusOK, "admin_posts.tmpl", h.addCommonData(c, gin.H{
		"title": "Posts",
		"posts": posts,
	}))
}

// ShowCreatePost displays the post creation form
func (h *AdminHandlers) ShowCreatePost(c *gin.Context) {
	tags, err := h.repos.Tags.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
		"title": "Create Post",
		"tags":  tags,
	}))
}

// CreatePost handles post creation
func (h *AdminHandlers) CreatePost(c *gin.Context) {
	// Get the logged in user
	userInterface, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "User session not found",
		}))
		return
	}
	user := userInterface.(*models.User)

	// Parse form data
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	excerpt := c.PostForm("excerpt")
	publishedAt := c.PostForm("publishedAt")
	var parsedTime time.Time

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to get settings",
		}))
		return
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	if publishedAt == "" {
		// If no date provided, use current time
		parsedTime = time.Now().In(loc)
	} else {
		var err error
		parsedTime, err = time.ParseInLocation("2006-01-02T15:04", publishedAt, loc)
		if err != nil {
			c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
				"error": "Invalid date format",
			}))
			return
		}
	}
	visible := c.PostForm("visible") == "on"

	// Basic validation
	if title == "" || slug == "" || content == "" {
		c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Title, slug and content are required",
		}))
		return
	}

	// Create post
	post := &models.Post{
		Title:       title,
		Slug:        slug,
		Content:     content,
		PublishedAt: parsedTime.UTC(),
		Visible:     visible,
		AuthorID:    user.ID,
	}

	// Only set excerpt if it's not empty
	if excerpt != "" {
		post.Excerpt = &excerpt
	}

	// Handle tags
	var tags []string
	if err := c.Request.ParseForm(); err != nil {
		c.HTML(http.StatusBadRequest, "admin_error.tmpl", gin.H{
			"error": "Invalid form data",
		})
		return
	}

	tags = strings.Split(c.PostForm("tags"), ",")

	// Create post with transaction to ensure atomic operation
	if err := h.repos.Posts.Create(post); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to create post",
			"post":  post,
		}))
		return
	}

	if err := h.repos.Posts.AssociateTags(post, tags); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to associate tags",
			"post":  post,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/posts")
}

// ListPostsByTag shows all posts for a specific tag
func (h *AdminHandlers) ListPostsByTag(c *gin.Context) {
	id := c.Param("id")

	tagID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Find the tag
	tag, err := h.tagRepo.FindByID(tagID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	posts, err := h.repos.Posts.FindByTag(tag.Slug)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
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

	c.HTML(http.StatusOK, "admin_tag_posts.tmpl", h.addCommonData(c, gin.H{
		"title": fmt.Sprintf("Posts tagged with '%s'", tag.Name),
		"posts": posts,
		"tag":   tag,
	}))
}

// ConfirmDeletePost shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePost(c *gin.Context) {
	id := c.Param("id")
	tagID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	post, err := h.postRepo.FindByID(tagID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_confirm_delete_post.tmpl", h.addCommonData(c, gin.H{
		"title": "Confirm Delete Post",
		"post":  post,
	}))
}

// DeletePost removes a post and its tag associations
func (h *AdminHandlers) DeletePost(c *gin.Context) {
	id := c.Param("id")
	tagID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	post, err := h.postRepo.FindByID(tagID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	if err := h.repos.Posts.Delete(post); err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *AdminHandlers) EditPost(c *gin.Context) {
	var allTags []*models.Tag
	id := c.Param("id")
	tagID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	post, err := h.postRepo.FindByID(tagID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	allTags, err = h.tagRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	// Convert post time to user timezone
	post.PublishedAt = post.PublishedAt.In(loc)

	c.HTML(http.StatusOK, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
		"title":    "Edit Post",
		"post":     post,
		"allTags":  allTags,
		"postTags": post.Tags,
	}))
}

func (h *AdminHandlers) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	tagID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	post, err := h.postRepo.FindByID(tagID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Parse form data
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	excerpt := c.PostForm("excerpt")
	publishedAt := c.PostForm("publishedAt")
	visible := c.PostForm("visible") == "on"

	// Basic validation
	if title == "" || slug == "" || content == "" {
		c.HTML(http.StatusBadRequest, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
			"error": "All fields are required",
			"post":  post,
		}))
		return
	}

	// Get settings for timezone
	settings, err := h.repos.Settings.Get()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to get settings",
			"post":  post,
		}))
		return
	}

	// Load timezone from settings
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.UTC
	}

	var parsedTime time.Time
	if publishedAt == "" {
		parsedTime = time.Now().In(loc)
	} else {
		var err error
		parsedTime, err = time.ParseInLocation("2006-01-02T15:04", publishedAt, loc)
		if err != nil {
			c.HTML(http.StatusBadRequest, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
				"error": "Invalid date format",
				"post":  post,
			}))
			return
		}
	}

	// Update post fields
	post.Title = title
	post.Slug = slug
	post.Content = content
	post.PublishedAt = parsedTime.UTC()
	post.Visible = visible

	// Update excerpt - set to nil if empty, otherwise update the value
	if excerpt == "" {
		post.Excerpt = nil
	} else {
		post.Excerpt = &excerpt
	}

	// Handle tags
	var tags []string
	if err := c.Request.ParseForm(); err != nil {
		c.HTML(http.StatusBadRequest, "admin_error.tmpl", gin.H{
			"error": "Invalid form data",
		})
		return
	}

	tags = strings.Split(c.PostForm("tags"), ",")

	// Update post with transaction to ensure atomic operation
	if err := h.repos.Posts.Update(post); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update post",
			"post":  post,
		}))
		return
	}

	if err := h.repos.Posts.AssociateTags(post, tags); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update tags",
			"post":  post,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/posts")
}
