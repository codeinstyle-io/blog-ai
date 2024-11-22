package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"codeinstyle.io/blog/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandlers struct {
	db *gorm.DB
}

func NewAdminHandlers(database *gorm.DB) *AdminHandlers {
	return &AdminHandlers{db: database}
}

// ListTags shows all tags and their post counts
func (h *AdminHandlers) ListTags(c *gin.Context) {
	var tags []struct {
		db.Tag
		PostCount int64
	}

	result := h.db.Model(&db.Tag{}).
		Select("tags.*, count(post_tags.post_id) as post_count").
		Joins("left join post_tags on post_tags.tag_id = tags.id").
		Group("tags.id").
		Find(&tags)

	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "errors/500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_tags.tmpl", gin.H{
		"title": "Tags",
		"tags":  tags,
	})
}

// DeleteTag removes a tag without affecting posts
func (h *AdminHandlers) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&db.Tag{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}

// ListUsers shows all users (except sensitive data)
func (h *AdminHandlers) ListUsers(c *gin.Context) {
	var users []db.User
	if err := h.db.Select("id, first_name, last_name, email, created_at, updated_at").Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_users.tmpl", gin.H{
		"title": "Users",
		"users": users,
	})
}

// ShowCreatePost displays the post creation form
func (h *AdminHandlers) ShowCreatePost(c *gin.Context) {
	var tags []db.Tag
	if err := h.db.Find(&tags).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_create_post.tmpl", gin.H{
		"title": "Create Post",
		"tags":  tags,
	})
}

func (h *AdminHandlers) CreatePost(c *gin.Context) {
	var post db.Post

	// Parse form data
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	publishedAt := c.PostForm("publishedAt")
	visible := c.PostForm("visible") == "on"

	// Basic validation
	if title == "" || slug == "" || content == "" || publishedAt == "" {
		c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", gin.H{
			"error": "All fields are required",
		})
		return
	}

	// Parse the published date
	parsedTime, err := time.Parse("2006-01-02T15:04", publishedAt)
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", gin.H{
			"error": "Invalid date format",
		})
		return
	}

	// Create post object
	post = db.Post{
		Title:       title,
		Slug:        slug,
		Content:     content,
		PublishedAt: parsedTime,
		Visible:     visible,
	}

	// Handle tags
	var tagNames []string
	if err := json.Unmarshal([]byte(c.PostForm("tags")), &tagNames); err != nil {
		c.HTML(http.StatusBadRequest, "admin/create_post.tmpl", gin.H{
			"error": "Invalid tags format",
			"post":  post,
		})
		return
	}

	// Create/get tags and associate
	var tags []db.Tag
	for _, name := range tagNames {
		var tag db.Tag
		result := h.db.Where(db.Tag{Name: name}).FirstOrCreate(&tag)
		if result.Error != nil {
			c.HTML(http.StatusInternalServerError, "admin/create_post.tmpl", gin.H{
				"error": "Failed to create tag",
				"post":  post,
			})
			return
		}
		tags = append(tags, tag)
	}
	post.Tags = tags

	// Create post with transaction to ensure atomic operation
	tx := h.db.Begin()
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin/create_post.tmpl", gin.H{
			"error": "Failed to create post",
			"post":  post,
		})
		return
	}

	if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin/create_post.tmpl", gin.H{
			"error": "Failed to associate tags",
			"post":  post,
		})
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}

// ListPosts shows all posts for admin
func (h *AdminHandlers) ListPosts(c *gin.Context) {
	var posts []db.Post
	if err := h.db.Preload("Tags").Find(&posts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_posts.tmpl", gin.H{
		"title": "Posts",
		"posts": posts,
	})
}

// ConfirmDeletePost shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post
	if err := h.db.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "errors/404.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_confirm_delete.tmpl", gin.H{
		"title": "Confirm Delete",
		"post":  post,
	})
}

// DeletePost removes a post and its tag associations
func (h *AdminHandlers) DeletePost(c *gin.Context) {
	id := c.Param("id")

	// Start transaction
	tx := h.db.Begin()

	// Delete post_tags associations
	if err := tx.Exec("DELETE FROM post_tags WHERE post_id = ?", id).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "errors/500.tmpl", gin.H{
			"error": "Failed to delete post tags",
		})
		return
	}

	// Delete post
	if err := tx.Delete(&db.Post{}, id).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "errors/500.tmpl", gin.H{
			"error": "Failed to delete post",
		})
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}

func (h *AdminHandlers) EditPost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post

	if err := h.db.Preload("Tags").First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "errors/404.tmpl", gin.H{})
		return
	}

	var allTags []db.Tag
	if err := h.db.Find(&allTags).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_edit_post.tmpl", gin.H{
		"title":   "Edit Post",
		"post":    post,
		"allTags": allTags,
	})
}

func (h *AdminHandlers) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post

	if err := h.db.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "errors/404.tmpl", gin.H{})
		return
	}

	// Update fields
	post.Title = c.PostForm("title")
	post.Slug = c.PostForm("slug")
	post.Content = c.PostForm("content")
	post.Visible = c.PostForm("visible") == "on"

	publishedAt, err := time.Parse("2006-01-02T15:04", c.PostForm("publishedAt"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_edit_post.tmpl", gin.H{
			"error": "Invalid date format",
			"post":  post,
		})
		return
	}
	post.PublishedAt = publishedAt

	// Handle tags
	var tagNames []string
	if err := json.Unmarshal([]byte(c.PostForm("tags")), &tagNames); err != nil {
		c.HTML(http.StatusBadRequest, "admin_edit_post.tmpl", gin.H{
			"error": "Invalid tags format",
			"post":  post,
		})
		return
	}

	// Update tags
	var tags []db.Tag
	for _, name := range tagNames {
		var tag db.Tag
		h.db.FirstOrCreate(&tag, db.Tag{Name: name})
		tags = append(tags, tag)
	}
	post.Tags = tags

	// Update with transaction
	tx := h.db.Begin()
	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin/edit_post.tmpl", gin.H{
			"error": "Failed to update post",
			"post":  post,
		})
		return
	}

	if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin/edit_post.tmpl", gin.H{
			"error": "Failed to update tags",
			"post":  post,
		})
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}

func (h *AdminHandlers) Index(c *gin.Context) {
	var postCount, tagCount, userCount int64
	var recentPosts []db.Post

	// Get counts
	h.db.Model(&db.Post{}).Count(&postCount)
	h.db.Model(&db.Tag{}).Count(&tagCount)
	h.db.Model(&db.User{}).Count(&userCount)

	// Get 5 most recent posts
	h.db.Order("published_at desc").Limit(5).Find(&recentPosts)

	c.HTML(http.StatusOK, "admin_index.tmpl", gin.H{
		"title":       "Dashboard",
		"postCount":   postCount,
		"tagCount":    tagCount,
		"userCount":   userCount,
		"recentPosts": recentPosts,
	})
}

// Add response struct
type tagResponse struct {
	Id   uint   `json:"id"`   // lowercase for JS
	Name string `json:"name"` // lowercase for JS
}

func (h *AdminHandlers) GetTags(c *gin.Context) {
	var dbTags []db.Tag

	if err := h.db.Find(&dbTags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	// Transform to response format
	tags := make([]tagResponse, len(dbTags))
	for i, tag := range dbTags {
		tags[i] = tagResponse{
			Id:   tag.ID,
			Name: tag.Name,
		}
	}

	c.JSON(http.StatusOK, tags)
}
