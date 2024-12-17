package handlers

import (
	"net/http"

	"codeinstyle.io/captain/cmd"
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// AuthHandlers handles all authentication related routes
type AuthHandlers struct {
	*BaseHandlers
	sessionStore *session.Store
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(repos *repository.Repositories, cfg *config.Config, sessionStore *session.Store) *AuthHandlers {
	return &AuthHandlers{
		BaseHandlers: NewBaseHandlers(repos, cfg),
		sessionStore: sessionStore,
	}
}

func (h *AuthHandlers) ShowLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func (h *AuthHandlers) PostLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	next := c.FormValue("next")
	if next == "" {
		next = "/admin"
	}

	if err := cmd.ValidateEmail(email); err != nil {
		return c.Status(http.StatusBadRequest).Render("login", fiber.Map{
			"error": "Invalid form data",
		})
	}

	// Find user by email
	user, err := h.repos.Users.FindByEmail(email)
	if err != nil {
		// Timing attack prevention
		utils.CheckPasswordHash(password, "")
		return c.Status(http.StatusUnauthorized).Render("login", fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check password
	if !utils.CheckPasswordHash(password, user.Password) {
		return c.Status(http.StatusUnauthorized).Render("login", fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Set session
	sess, err := h.sessionStore.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("login", fiber.Map{
			"error": "Failed to create session",
		})
	}

	sess.Set("loggedIn", true)

	if err := sess.Save(); err != nil {
		return c.Status(http.StatusInternalServerError).Render("login", fiber.Map{
			"error": "Failed to save session",
		})
	}

	return c.Redirect(next)
}

func (h *AuthHandlers) Logout(c *fiber.Ctx) error {
	sess, err := h.sessionStore.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{})
	}

	return c.Redirect("/login")
}

// HandleSetup handles both GET and POST requests for the setup page
func (h *AuthHandlers) HandleSetup(c *fiber.Ctx) error {
	// If users already exist, redirect to login
	count, err := h.repos.Users.CountAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("setup", fiber.Map{"Error": "Failed to count users"})
	}

	if count > 0 {
		return c.Redirect("/login?next=/admin")
	}

	// Handle POST request
	if c.Method() == fiber.MethodPost {
		email := c.FormValue("email")
		password := c.FormValue("password")
		firstName := c.FormValue("firstName")
		lastName := c.FormValue("lastName")

		// Validate input
		if err := cmd.ValidateEmail(email); err != nil {
			return c.Status(http.StatusBadRequest).Render("setup", fiber.Map{"Error": "Invalid email address"})
		}
		if err := cmd.ValidatePassword(password); err != nil {
			return c.Status(http.StatusBadRequest).Render("setup", fiber.Map{"Error": "Password must be at least 8 characters"})
		}
		if err := cmd.ValidateFirstName(firstName); err != nil {
			return c.Status(http.StatusBadRequest).Render("setup", fiber.Map{"Error": err.Error()})
		}
		if err := cmd.ValidateLastName(lastName); err != nil {
			return c.Status(http.StatusBadRequest).Render("setup", fiber.Map{"Error": err.Error()})
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return c.Status(http.StatusInternalServerError).Render("setup", fiber.Map{"Error": "Failed to hash password"})
		}

		// Create admin user
		user := &models.User{
			Email:     email,
			Password:  hashedPassword,
			FirstName: firstName,
			LastName:  lastName,
		}

		if err := h.repos.Users.Create(user); err != nil {
			return c.Status(http.StatusInternalServerError).Render("setup", fiber.Map{"Error": "Failed to create user"})
		}

		// Redirect to admin login
		return c.Redirect("/login")
	}

	// Handle GET request
	return c.Render("setup", fiber.Map{})
}
