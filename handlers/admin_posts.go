package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
)

// ListPosts shows all posts for admin
func (h *AdminHandlers) ListPosts(c *gin.Context) {
	var posts []db.Post
	if err := h.db.Preload("Tags").Preload("Author").Find(&posts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_posts.tmpl", h.addCommonData(c, gin.H{
		"title": "Posts",
		"posts": posts,
	}))
}

// ShowCreatePost displays the post creation form
func (h *AdminHandlers) ShowCreatePost(c *gin.Context) {
	var tags []db.Tag
	if err := h.db.Find(&tags).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
		"title": "Create Post",
		"tags":  tags,
	}))
}

func (h *AdminHandlers) CreatePost(c *gin.Context) {
	// Get the logged in user
	userInterface, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "User session not found",
		}))
		return
	}
	user := userInterface.(*db.User)

	var post db.Post

	// Parse form data
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	excerpt := c.PostForm("excerpt")
	publishedAt := c.PostForm("publishedAt")
	var parsedTime time.Time

	// Get settings for timezone
	settings, err := db.GetSettings(h.db)
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
	post = db.Post{
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
	var tagNames []string
	tagsJSON := c.PostForm("tags")
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &tagNames); err != nil {
			c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
				"error": "Invalid tags format",
				"post":  post,
			}))
			return
		}
	}

	// Create/get tags and associate
	var tags []db.Tag
	for _, name := range tagNames {
		var tag db.Tag
		result := h.db.Where(db.Tag{Name: name}).FirstOrCreate(&tag)
		if result.Error != nil {
			c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
				"error": "Failed to create tag",
				"post":  post,
			}))
			return
		}
		tags = append(tags, tag)
	}
	post.Tags = tags

	// Create post with transaction to ensure atomic operation
	tx := h.db.Begin()
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to create post",
			"post":  post,
		}))
		return
	}

	if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to associate tags",
			"post":  post,
		}))
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}

// ListPostsByTag shows all posts for a specific tag
func (h *AdminHandlers) ListPostsByTag(c *gin.Context) {
	tagID := c.Param("id")
	var tag db.Tag
	if err := h.db.First(&tag, tagID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	var posts []db.Post
	if err := h.db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.id = ?", tagID).
		Preload("Tags").
		Preload("Author").
		Find(&posts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_tag_posts.tmpl", h.addCommonData(c, gin.H{
		"title": "Posts",
		"posts": posts,
		"tag":   tag,
	}))
}

// ConfirmDeletePost shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post
	if err := h.db.First(&post, id).Error; err != nil {
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

	// Start transaction
	tx := h.db.Begin()

	// First, find the post
	var post db.Post
	if err := tx.First(&post, id).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	fmt.Println("Deleting post with ID:", id)

	// Clear associations first
	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		fmt.Println("Error clearing associations:", err)
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Then delete the post
	if err := tx.Delete(&post).Error; err != nil {
		fmt.Println("Error deleting post:", err)
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *AdminHandlers) EditPost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post
	if err := h.db.Preload("Tags").First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	var allTags []db.Tag
	if err := h.db.Find(&allTags).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
		"title":    "Edit Post",
		"post":     post,
		"allTags":  allTags,
		"postTags": post.Tags,
	}))
}

func (h *AdminHandlers) UpdatePost(c *gin.Context) {
	// Get the post ID from the URL
	postID := c.Param("id")

	// Find the existing post
	var post db.Post
	if err := h.db.First(&post, postID).Error; err != nil {
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
	settings, err := db.GetSettings(h.db)
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
	var tagNames []string
	tagsJSON := c.PostForm("tags")
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &tagNames); err != nil {
			c.HTML(http.StatusBadRequest, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
				"error": "Invalid tags format",
				"post":  post,
			}))
			return
		}
	}

	// Create/get tags
	var tags []db.Tag
	for _, name := range tagNames {
		var tag db.Tag
		result := h.db.Where(db.Tag{Name: name}).FirstOrCreate(&tag)
		if result.Error != nil {
			c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
				"error": "Failed to create tag",
				"post":  post,
			}))
			return
		}
		tags = append(tags, tag)
	}

	// Start transaction for update
	tx := h.db.Begin()

	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update post",
			"post":  post,
		}))
		return
	}

	if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update tags",
			"post":  post,
		}))
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}
