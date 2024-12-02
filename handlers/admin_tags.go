package handlers

import (
	"net/http"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
)

// tagResponse struct for API responses
type tagResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
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
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
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

// ShowCreateTag displays the tag creation form
func (h *AdminHandlers) ShowCreateTag(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_tag.tmpl", gin.H{
		"title": "Create Tag",
	})
}

// CreateTag handles tag creation
func (h *AdminHandlers) CreateTag(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.HTML(http.StatusBadRequest, "admin_create_tag.tmpl", gin.H{
			"error": "Tag name is required",
		})
		return
	}

	tag := db.Tag{
		Name: name,
	}

	if err := h.db.Create(&tag).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: tags.name" {
			c.HTML(http.StatusBadRequest, "admin_create_tag.tmpl", gin.H{
				"error": "Tag name already exists",
			})
			return
		}

		c.HTML(http.StatusInternalServerError, "admin_create_tag.tmpl", gin.H{
			"error": "Failed to create tag",
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/tags")
}

// GetTags returns a list of tags for API consumption
func (h *AdminHandlers) GetTags(c *gin.Context) {
	var tags []db.Tag
	if err := h.db.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	var response []tagResponse
	for _, tag := range tags {
		response = append(response, tagResponse{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}
