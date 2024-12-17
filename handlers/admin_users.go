package handlers

import (
	"net/http"

	"codeinstyle.io/captain/cmd"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ListUsers shows all users (except sensitive data)
func (h *AdminHandlers) ListUsers(c *fiber.Ctx) error {
	users, err := h.repos.Users.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{})
	}

	return c.Render("admin_users", fiber.Map{
		"users": users,
	})
}

// ShowCreateUser displays the user creation form
func (h *AdminHandlers) ShowCreateUser(c *fiber.Ctx) error {
	return c.Render("admin_create_user", fiber.Map{
		"user": &models.User{},
	})
}

// CreateUser handles user creation
func (h *AdminHandlers) CreateUser(c *fiber.Ctx) error {
	firstName := c.FormValue("firstName")
	lastName := c.FormValue("lastName")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Validate input
	if err := cmd.ValidateFirstName(firstName); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"user":  &models.User{},
			"error": err.Error(),
		})
	}
	if err := cmd.ValidateLastName(lastName); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"user":  &models.User{},
			"error": err.Error(),
		})
	}
	if err := cmd.ValidateEmail(email); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"user":  &models.User{},
			"error": err.Error(),
		})
	}
	if err := cmd.ValidatePassword(password); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"user":  &models.User{},
			"error": err.Error(),
		})
	}

	// Check if email already exists
	count, err := h.repos.Users.CountByEmail(email)

	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_create_user", fiber.Map{
			"user":  &models.User{},
			"error": "Failed to check email uniqueness",
		})
	}
	if count > 0 {
		return c.Status(http.StatusBadRequest).Render("admin_create_user", fiber.Map{
			"user":  &models.User{},
			"error": "Email already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_create_user", fiber.Map{
			"user":  &models.User{},
			"error": "Failed to hash password",
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
		return c.Status(http.StatusInternalServerError).Render("admin_create_user", fiber.Map{
			"user":  user,
			"error": "Failed to create user",
		})
	}

	return c.Redirect("/admin/users")
}

// ShowEditUser displays the user edit form
func (h *AdminHandlers) ShowEditUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("admin_edit_user", fiber.Map{
		"user": user,
	})
}

// UpdateUser handles user updates
func (h *AdminHandlers) UpdateUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"user":  user,
			"error": "Invalid form data",
		})
	}

	// Validate input
	if err := cmd.ValidateFirstName(user.FirstName); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"user":  user,
			"error": err.Error(),
		})
	}
	if err := cmd.ValidateLastName(user.LastName); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"user":  user,
			"error": err.Error(),
		})
	}
	if err := cmd.ValidateEmail(user.Email); err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"user":  user,
			"error": err.Error(),
		})
	}

	// Check if email already exists for other users
	count, err := h.repos.Users.CountByEmail(user.Email)

	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_edit_user", fiber.Map{
			"user":  user,
			"error": "Failed to check email uniqueness",
		})
	}

	if count > 1 {
		return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
			"user":  user,
			"error": "Email already exists",
		})
	}

	// Update password if provided
	if user.Password != "" {
		if err := cmd.ValidatePassword(user.Password); err != nil {
			return c.Status(http.StatusBadRequest).Render("admin_edit_user", fiber.Map{
				"user":  user,
				"error": err.Error(),
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Render("admin_edit_user", fiber.Map{
				"user":  user,
				"error": "Failed to hash password",
			})
		}
		user.Password = string(hashedPassword)
	}

	// Update user
	if err := h.repos.Users.Update(user); err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_edit_user", fiber.Map{
			"user":  user,
			"error": "Failed to update user",
		})
	}

	return c.Redirect("/admin/users")
}

// ShowDeleteUser displays the user deletion confirmation page
func (h *AdminHandlers) ShowDeleteUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_user", fiber.Map{
		"user": user,
	})
}

// DeleteUser handles user deletion
func (h *AdminHandlers) DeleteUser(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("500", fiber.Map{})
	}

	user, err := h.repos.Users.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	// Delete user
	if err := h.repos.Users.Delete(user); err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{})
	}

	return c.Redirect("/admin/users")
}
