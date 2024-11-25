package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandlers struct {
	db *gorm.DB
}

func NewPostHandlers(database *gorm.DB) *PostHandlers {
	return &PostHandlers{db: database}
}

func (h *PostHandlers) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var post db.Post
	if err := h.db.Where("slug = ?", slug).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
				"title": "Post not found",
			})
			return
		}
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
			"title": "Error",
		})
		return
	}

	c.HTML(http.StatusOK, "post.tmpl", gin.H{
		"title": post.Title,
		"post":  post,
	})
}

func (h *PostHandlers) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 3 // Posts per page

	var total int64
	h.db.Model(&db.Post{}).Where("visible = ?", true).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	offset := (page - 1) * perPage

	var posts []db.Post
	result := h.db.Preload("Tags").
		Where("visible = ?", true).
		Order("published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts)

	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "posts.tmpl", gin.H{
		"title":       "Latest Articles",
		"posts":       posts,
		"currentPage": page,
		"totalPages":  totalPages,
	})
}

func (h *PostHandlers) ListPostsByTag(c *gin.Context) {
	tagName := c.Param("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 3 // Posts per page

	// Get total count
	var total int64
	h.db.Model(&db.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.name = ? AND posts.visible = ?", tagName, true).
		Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	offset := (page - 1) * perPage

	// Get posts
	var posts []db.Post
	result := h.db.Preload("Tags").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.name = ? AND posts.visible = ?", tagName, true).
		Order("published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts)

	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "tag_posts.tmpl", gin.H{
		"title":       fmt.Sprintf("Posts tagged with #%s", tagName),
		"tag":         tagName,
		"posts":       posts,
		"currentPage": page,
		"totalPages":  totalPages,
	})
}

// ...other post handlers like GetPostBySlug, EditPost, etc...
