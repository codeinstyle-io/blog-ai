package handlers

import (
	"net/http"

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
		"title":    "Login",
		"returnTo": c.Query("returnTo"),
	}))
}

func (h *AuthHandlers) PostLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	returnTo := c.Query("returnTo")
	if returnTo == "" {
		returnTo = "/admin"
	}

	user, err := db.GetUserByEmail(h.db, email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.tmpl", h.addCommonData(gin.H{
			"title":    "Login",
			"error":    "Invalid credentials",
			"returnTo": returnTo,
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
	c.Redirect(http.StatusFound, returnTo)
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
