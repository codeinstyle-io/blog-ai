package handlers

import (
	"net/http"

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
