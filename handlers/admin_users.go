package handlers

import (
	"net/http"

	"codeinstyle.io/captain/cmd"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
)

// ListUsers shows all users (except sensitive data)
func (h *AdminHandlers) ListUsers(c *gin.Context) {

	users, err := h.userRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{
			"error": err.Error(),
		}))
		return
	}

	c.HTML(http.StatusOK, "admin_users.tmpl", h.addCommonData(c, gin.H{
		"title": "Users",
		"users": users,
	}))
}

// ShowCreateUser displays the user creation form
func (h *AdminHandlers) ShowCreateUser(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
		"title": "Create User",
	}))
}

// CreateUser handles user creation
func (h *AdminHandlers) CreateUser(c *gin.Context) {
	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate input
	if err := cmd.ValidateFirstName(firstName); err != nil {
		c.HTML(http.StatusBadRequest, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": err.Error(),
		}))
		return
	}
	if err := cmd.ValidateLastName(lastName); err != nil {
		c.HTML(http.StatusBadRequest, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": err.Error(),
		}))
		return
	}
	if err := cmd.ValidateEmail(email); err != nil {
		c.HTML(http.StatusBadRequest, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": err.Error(),
		}))
		return
	}
	if err := cmd.ValidatePassword(password); err != nil {
		c.HTML(http.StatusBadRequest, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": err.Error(),
		}))
		return
	}

	// Check if email already exists
	count, err := h.userRepo.CountByEmail(email)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": "Failed to check email uniqueness",
		}))
		return
	}
	if count > 0 {
		c.HTML(http.StatusBadRequest, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": "Email already exists",
		}))
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": "Failed to hash password",
		}))
		return
	}

	// Create user
	user := &models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPassword,
	}

	if err := h.userRepo.Create(user); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Create User",
			"error": "Failed to create user",
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// ShowEditUser displays the user edit form
func (h *AdminHandlers) ShowEditUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	c.HTML(http.StatusOK, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
		"title": "Edit User",
		"user":  user,
	}))
}

// UpdateUser handles user updates
func (h *AdminHandlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate input
	if err := cmd.ValidateFirstName(firstName); err != nil {
		c.HTML(http.StatusBadRequest, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Edit User",
			"user":  user,
			"error": err.Error(),
		}))
		return
	}
	if err := cmd.ValidateLastName(lastName); err != nil {
		c.HTML(http.StatusBadRequest, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Edit User",
			"user":  user,
			"error": err.Error(),
		}))
		return
	}
	if err := cmd.ValidateEmail(email); err != nil {
		c.HTML(http.StatusBadRequest, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Edit User",
			"user":  user,
			"error": err.Error(),
		}))
		return
	}

	// Check if email already exists for other users
	count, err := h.userRepo.CountByEmail(email)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Edit User",
			"user":  user,
			"error": "Failed to check email uniqueness",
		}))
		return
	}

	if count > 1 {
		c.HTML(http.StatusBadRequest, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Edit User",
			"user":  user,
			"error": "Email already exists",
		}))
		return
	}

	// Update user fields
	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email

	// Update password if provided
	if password != "" {
		if err := cmd.ValidatePassword(password); err != nil {
			c.HTML(http.StatusBadRequest, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
				"title": "Edit User",
				"user":  user,
				"error": err.Error(),
			}))
			return
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
				"title": "Edit User",
				"user":  user,
				"error": "Failed to hash password",
			}))
			return
		}
		user.Password = hashedPassword
	}

	if err := h.userRepo.Update(user); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_user.tmpl", h.addCommonData(c, gin.H{
			"title": "Edit User",
			"user":  user,
			"error": "Failed to update user",
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// ShowDeleteUser displays the user deletion confirmation page
func (h *AdminHandlers) ShowDeleteUser(c *gin.Context) {
	var postCount int64
	id := c.Param("id")
	userID, err := utils.ParseUint(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

	// Check if user has any posts
	postCount, err = h.postRepo.CountByAuthor(user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{
			"error": err.Error(),
		}))
		return
	}

	c.HTML(http.StatusOK, "admin_confirm_delete_user.tmpl", h.addCommonData(c, gin.H{
		"title":      "Delete User",
		"user":       user,
		"hasContent": postCount > 0,
	}))
}

// DeleteUser handles user deletion
func (h *AdminHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := utils.ParseUint(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete user (posts will remain with the author info)
	if err := h.userRepo.Delete(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
