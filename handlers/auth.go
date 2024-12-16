package handlers

import (
	"net/http"
	"time"

	"codeinstyle.io/captain/cmd"
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
)

// AuthHandlers handles all authentication related routes
type AuthHandlers struct {
	*BaseHandler
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(repos *repository.Repositories, cfg *config.Config) *AuthHandlers {
	return &AuthHandlers{
		BaseHandler: NewBaseHandler(repos, cfg),
	}
}

func (h *AuthHandlers) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", h.addCommonData(gin.H{
		"title": "Login",
	}))
}

func (h *AuthHandlers) PostLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	next := c.PostForm("next")
	if next == "" {
		next = "/admin"
	}

	// Check if user exists
	user, _ := h.userRepo.FindByEmail(email)

	if user == nil {
		user = &models.User{
			Password: "",
		}
	}

	// Compare hashed password
	if !utils.CheckPasswordHash(password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.tmpl", h.addCommonData(gin.H{
			"title": "Login",
			"error": "Invalid credentials",
			"next":  next,
		}))
		return
	}

	// Generate session token
	token, err := utils.GenerateSessionToken()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
			"title": "Login",
			"error": "Failed to generate session token",
		}))
		return
	}

	// Create session
	session := &models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := h.sessionRepo.Create(session); err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
			"title": "Login",
			"error": "Failed to create session",
		}))
		return
	}

	// Set cookie
	c.SetCookie(system.CookieName, token, int(24*time.Hour.Seconds()), "/", h.config.Site.Domain, h.config.Site.SecureCookie, true)
	c.Redirect(http.StatusFound, next)
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	token, err := c.Cookie(system.CookieName)
	if err == nil {
		// Delete session from database
		if err := h.sessionRepo.DeleteByToken(token); err != nil {
			c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
				"error": "Failed to delete session",
			}))
			return
		}
	}

	// Clear cookie
	c.SetCookie(system.CookieName, "", -1, "/", h.config.Site.Domain, h.config.Site.SecureCookie, true)
	c.Redirect(http.StatusFound, "/login")
}

// HandleSetup handles both GET and POST requests for the setup page
func (h *AuthHandlers) HandleSetup(c *gin.Context) {
	// If users already exist, redirect to login
	count, err := h.userRepo.CountAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "setup.tmpl", gin.H{"Error": "Failed to count users"})
		return
	}

	if count > 0 {
		c.Redirect(http.StatusFound, "/login?next=/admin")
		return
	}

	// Handle POST request
	if c.Request.Method == http.MethodPost {
		email := c.PostForm("email")
		password := c.PostForm("password")
		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")

		// Validate input
		if err := cmd.ValidateEmail(email); err != nil {
			c.HTML(http.StatusBadRequest, "setup.tmpl", gin.H{"Error": "Invalid email address"})
			return
		}
		if err := cmd.ValidatePassword(password); err != nil {
			c.HTML(http.StatusBadRequest, "setup.tmpl", gin.H{"Error": "Password must be at least 8 characters"})
			return
		}
		if err := cmd.ValidateFirstName(firstName); err != nil {
			c.HTML(http.StatusBadRequest, "setup.tmpl", gin.H{"Error": err.Error()})
			return
		}
		if err := cmd.ValidateLastName(lastName); err != nil {
			c.HTML(http.StatusBadRequest, "setup.tmpl", gin.H{"Error": err.Error()})
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "setup.tmpl", gin.H{"Error": "Failed to hash password"})
			return
		}

		// Create admin user
		user := &models.User{
			Email:     email,
			Password:  hashedPassword,
			FirstName: firstName,
			LastName:  lastName,
		}

		if err := h.userRepo.Create(user); err != nil {
			c.HTML(http.StatusInternalServerError, "setup.tmpl", gin.H{"Error": "Failed to create user"})
			return
		}

		// Redirect to admin login
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Handle GET request
	c.HTML(http.StatusOK, "setup.tmpl", gin.H{})
}

func (h *AuthHandlers) addCommonData(data gin.H) gin.H {
	// Get menu items
	menuItems, _ := h.BaseHandler.menuRepo.FindAll()
	settings, _ := h.BaseHandler.repos.Settings.Get()

	// Add menu items to the data
	data["menuItems"] = menuItems

	// Add site config from settings
	data["config"] = gin.H{
		"SiteTitle":    settings.Title,
		"SiteSubtitle": settings.Subtitle,
		"Theme":        settings.Theme,
	}

	return data
}
