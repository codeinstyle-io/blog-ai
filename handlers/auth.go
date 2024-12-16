package handlers

import (
	"net/http"
	"time"

	"codeinstyle.io/captain/cmd"
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandlers handles all authentication related routes
type AuthHandlers struct {
	*BaseHandler
	sessionRepo *repository.SessionRepository
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
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := h.userRepo.FindByUsername(username)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", h.addCommonData(gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		}))
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", h.addCommonData(gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		}))
		return
	}

	// Generate session token
	token, err := utils.GenerateSessionToken()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
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
			"error": "Failed to create session",
		}))
		return
	}

	// Set cookie
	c.SetCookie("session_token", token, int(24*time.Hour.Seconds()), "/", h.config.Site.Domain, h.config.Site.SecureCookie, true)
	c.Redirect(http.StatusFound, "/admin")
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	token, err := c.Cookie("session_token")
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
	c.SetCookie("session_token", "", -1, "/", h.config.Site.Domain, h.config.Site.SecureCookie, true)
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

	// Add menu items to the data
	data["menuItems"] = menuItems

	// Add site config from settings
	data["config"] = gin.H{
		"SiteTitle":    h.BaseHandler.settings.Title,
		"SiteSubtitle": h.BaseHandler.settings.Subtitle,
		"Theme":        h.BaseHandler.settings.Theme,
	}

	return data
}
