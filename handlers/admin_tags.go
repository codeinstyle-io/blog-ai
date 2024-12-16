package handlers

import (
	"net/http"

	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
)

// tagResponse struct for API responses
type tagResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// ListTags shows all tags and their post counts
func (h *AdminHandlers) ListTags(c *gin.Context) {
	tagsAndCount, err := h.tagRepo.FindPostsAndCount()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_tags.tmpl", h.addCommonData(c, gin.H{
		"title": "Tags",
		"tags":  tagsAndCount,
	}))
}

// DeleteTag removes a tag without affecting posts
func (h *AdminHandlers) DeleteTag(c *gin.Context) {
	id := c.Param("id")

	tagID, err := utils.ParseUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	tag, err := h.tagRepo.FindByID(tagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	if err := h.tagRepo.Delete(tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

// ConfirmDeleteTag shows deletion confirmation page for a tag
func (h *AdminHandlers) ConfirmDeleteTag(c *gin.Context) {
	id := c.Param("id")

	tagID, err := utils.ParseUint(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	tag, err := h.tagRepo.FindByID(tagID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_confirm_delete_tag.tmpl", h.addCommonData(c, gin.H{
		"title": "Confirm Delete Tag",
		"tag":   tag,
	}))
}

// ShowCreateTag displays the tag creation form
func (h *AdminHandlers) ShowCreateTag(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_tag.tmpl", h.addCommonData(c, gin.H{
		"title": "Create Tag",
	}))
}

// CreateTag handles tag creation
func (h *AdminHandlers) CreateTag(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.HTML(http.StatusBadRequest, "admin_create_tag.tmpl", h.addCommonData(c, gin.H{
			"error": "Tag name is required",
		}))
		return
	}

	tag := models.Tag{
		Name: name,
	}

	if err := h.tagRepo.Create(&tag); err != nil {
		if err.Error() == "UNIQUE constraint failed: tags.name" {
			c.HTML(http.StatusBadRequest, "admin_create_tag.tmpl", h.addCommonData(c, gin.H{
				"error": "Tag name already exists",
			}))
			return
		}

		c.HTML(http.StatusInternalServerError, "admin_create_tag.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to create tag",
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/tags")
}

// GetTags returns a list of tags for API consumption
func (h *AdminHandlers) GetTags(c *gin.Context) {
	tags, err := h.repos.Tags.FindAll()
	if err != nil {
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
