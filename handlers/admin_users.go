package handlers

import (
	"net/http"

	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ListUsers shows all users (except sensitive data)
func (h *AdminHandlers) ListUsers(c *fiber.Ctx) error {
	users, err := h.repos.Users.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Render("admin_users", fiber.Map{
		"title": "Users",
		"users": users,
	})
}

// ShowCreateUser displays the user creation form
func (h *AdminHandlers) ShowCreateUser(c *fiber.Ctx) error {
	return c.Render("admin_create_user", fiber.Map{
		"title": "Create User",
		"user":  &models.User{},
	})
}

// CreateUser handles user creation
func (h *AdminHandlers) CreateUser(c *fiber.Ctx) error {
	firstName := c.FormValue("firstName")
	lastName := c.FormValue("lastName")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Validate input
	if err := utils.ValidateFirstName(firstName); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  &models.User{},
		})
	}
	if err := utils.ValidateLastName(lastName); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  &models.User{},
		})
	}
	if err := utils.ValidateEmail(email); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  &models.User{},
		})
	}
	if err := utils.ValidatePassword(password); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  &models.User{},
		})
	}

	// Check if email already exists
	count, err := h.repos.Users.CountByEmail(email)

	if err != nil {
		flash.Error(c, "Failed to check email uniqueness")
		return c.Status(http.StatusInternalServerError).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  &models.User{},
		})
	}
	if count > 0 {
		flash.Error(c, "Email already exists")
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  &models.User{},
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		flash.Error(c, "Failed to hash password")
		return c.Status(http.StatusInternalServerError).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  &models.User{},
		})
	}

	// Create user
	user := &models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
	}

	if err := h.repos.Users.Create(user); err != nil {
		flash.Error(c, "Failed to create user")
		return c.Status(http.StatusInternalServerError).Render("admin_create_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}

	flash.Success(c, "User created successfully")
	return c.Redirect("/admin/users")
}

// ShowEditUser displays the user edit form
func (h *AdminHandlers) ShowEditUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	return c.Render("admin_edit_user", fiber.Map{
		"title": "Edit User",
		"user":  user,
	})
}

// UpdateUser handles user updates
func (h *AdminHandlers) UpdateUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid user ID")
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		flash.Error(c, "User not found")
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	if err := c.BodyParser(&user); err != nil {
		flash.Error(c, "Invalid form data")
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}

	// Validate input
	if err := utils.ValidateFirstName(user.FirstName); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}
	if err := utils.ValidateLastName(user.LastName); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}
	if err := utils.ValidateEmail(user.Email); err != nil {
		flash.Error(c, err.Error())
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}

	// Check if email already exists for other users
	count, err := h.repos.Users.CountByEmail(user.Email)

	if err != nil {
		flash.Error(c, "Failed to check email uniqueness")
		return c.Status(http.StatusInternalServerError).Render("admin_edit_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}

	if count > 1 {
		flash.Error(c, "Email already exists")
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}

	// Update password if provided
	if user.Password != "" {
		if err := utils.ValidatePassword(user.Password); err != nil {
			flash.Error(c, err.Error())
			return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
				"title": "Users",
				"user":  user,
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			flash.Error(c, "Failed to hash password")
			return c.Status(http.StatusInternalServerError).Render("admin_edit_user", fiber.Map{
				"title": "Users",
				"user":  user,
			})
		}
		user.Password = string(hashedPassword)
	}

	// Update user
	if err := h.repos.Users.Update(user); err != nil {
		flash.Error(c, "Failed to update user")
		return c.Status(http.StatusInternalServerError).Render("admin_edit_user", fiber.Map{
			"title": "Users",
			"user":  user,
		})
	}

	flash.Success(c, "User updated successfully")
	return c.Redirect("/admin/users")
}

// ConfirmDeleteUser displays the user deletion confirmation page
func (h *AdminHandlers) ConfirmDeleteUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_user", fiber.Map{
		"title": "Confirm user deletion",
		"user":  user,
	})
}

// DeleteUser handles user deletion
func (h *AdminHandlers) DeleteUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid user ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid user ID",
			"redirect": "/admin/users",
		})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		flash.Error(c, "User not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":    "User not found",
			"redirect": "/admin/users",
		})
	}

	// Delete user
	if err := h.repos.Users.Delete(user); err != nil {
		flash.Error(c, "Failed to delete user")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to delete user",
			"redirect": "/admin/users",
		})
	}

	flash.Success(c, "User deleted successfully")
	return c.JSON(fiber.Map{
		"message":  "User deleted successfully",
		"redirect": "/admin/users",
	})
}
