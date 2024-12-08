package handlers

import (
	"net/http"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
)

// ListUsers shows all users (except sensitive data)
func (h *AdminHandlers) ListUsers(c *gin.Context) {
	var users []db.User
	if err := h.db.Select("id, first_name, last_name, email, created_at, updated_at").Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}
	c.HTML(http.StatusOK, "admin_users.tmpl", h.addCommonData(c, gin.H{
		"title": "Users",
		"users": users,
	}))
}
