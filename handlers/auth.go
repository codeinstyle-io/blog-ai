package handlers

import (
	"net/http"

	"codeinstyle.io/captain/cmd"
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandlers struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthHandlers(database *gorm.DB, config *config.Config) *AuthHandlers {
	return &AuthHandlers{db: database, config: config}
}

func (h *AuthHandlers) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", h.addCommonData(gin.H{
		"title": "Login",
		"next":  c.Query("next"),
	}))
}

func (h *AuthHandlers) PostLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	next := c.PostForm("next")
	if next == "" {
		next = "/admin"
	}

	user, err := db.GetUserByEmail(h.db, email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.tmpl", h.addCommonData(gin.H{
			"title": "Login",
			"error": "Invalid credentials",
			"next":  next,
		}))
		return
	}

	// Generate session token
	token, err := db.GenerateSessionToken()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
			"title": "Error",
		}))
		return
	}

	// Update user with session token
	user.SessionToken = &token
	if err := h.db.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
			"title": "Error",
		}))
		return
	}

	// Set session cookie
	c.SetCookie("session", token, 86400, "/", "", false, true)
	c.Redirect(http.StatusFound, next)
}

func (h *AuthHandlers) addCommonData(data gin.H) gin.H {
	// Get menu items
	var menuItems []db.MenuItem
	h.db.Preload("Page").Order("position").Find(&menuItems)

	// Add menu items to the data
	data["menuItems"] = menuItems

	// Add site config
	data["config"] = gin.H{
		"SiteTitle":    h.config.Site.Title,
		"SiteSubtitle": h.config.Site.Subtitle,
	}

	return data
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	// Clear session cookie
	c.SetCookie("session", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")

	// Save theme preference before logout
	theme, _ := c.Cookie("admin_theme")

	// Restore theme after redirect is set
	if theme != "" {
		c.SetCookie("admin_theme", theme, 3600*24*365, "/", "", false, false)
	}
}

// HandleSetup handles both GET and POST requests for the setup page
func (h *AuthHandlers) HandleSetup(c *gin.Context) {
	// If users already exist, redirect to login
	var count int64
	h.db.Model(&db.User{}).Count(&count)
	if count > 0 {
		c.Redirect(http.StatusFound, "/admin/login")
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
			c.HTML(http.StatusBadRequest, "pages/setup.tmpl", gin.H{"Error": "Invalid email address"})
			return
		}
		if err := cmd.ValidatePassword(password); err != nil {
			c.HTML(http.StatusBadRequest, "pages/setup.tmpl", gin.H{"Error": "Password must be at least 8 characters"})
			return
		}
		if err := cmd.ValidateFirstName(firstName); err != nil {
			c.HTML(http.StatusBadRequest, "pages/setup.tmpl", gin.H{"Error": err.Error()})
			return
		}
		if err := cmd.ValidateLastName(lastName); err != nil {
			c.HTML(http.StatusBadRequest, "pages/setup.tmpl", gin.H{"Error": err.Error()})
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "pages/setup.tmpl", gin.H{"Error": "Failed to hash password"})
			return
		}

		// Create admin user
		user := &db.User{
			Email:     email,
			Password:  hashedPassword,
			FirstName: firstName,
			LastName:  lastName,
		}

		if err := db.CreateUser(h.db, user); err != nil {
			c.HTML(http.StatusInternalServerError, "pages/setup.tmpl", gin.H{"Error": "Failed to create user"})
			return
		}

		// Redirect to admin login
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Handle GET request
	c.HTML(http.StatusOK, "setup.tmpl", gin.H{})
}
