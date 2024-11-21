package handlers

import (
	"net/http"
	"time"

	"codeinstyle.io/blog/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandlers struct {
	db *gorm.DB
}

func NewPostHandlers(database *gorm.DB) *PostHandlers {
	return &PostHandlers{db: database}
}

func (h *PostHandlers) CreatePost(c *gin.Context) {
	var post db.Post

	// Parse form data
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	publishedAt := c.PostForm("publishedAt")
	visible := c.PostForm("visible") == "on"

	// Basic validation
	if title == "" || slug == "" || content == "" || publishedAt == "" {
		c.HTML(http.StatusBadRequest, "create_post.tmpl", gin.H{
			"error": "All fields are required",
		})
		return
	}

	// Parse the published date
	parsedTime, err := time.Parse("2006-01-02T15:04", publishedAt)
	if err != nil {
		c.HTML(http.StatusBadRequest, "create_post.tmpl", gin.H{
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

	if err := h.db.Create(&post).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "create_post.tmpl", gin.H{
			"error": "Failed to create post",
		})
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}

func (h *PostHandlers) ListPosts(c *gin.Context) {
	posts, err := db.GetPosts(h.db, 5)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
			"title": "Error",
		})
		return
	}
	c.HTML(http.StatusOK, "posts.tmpl", gin.H{
		"title": "Latest articles",
		"posts": posts,
	})
}

// ...other post handlers like GetPostBySlug, EditPost, etc...
