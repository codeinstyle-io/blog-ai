package handlers

import (
	"net/http"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandlers struct {
	db *gorm.DB
}

func NewAuthHandlers(database *gorm.DB) *AuthHandlers {
	return &AuthHandlers{db: database}
}

func (h *AuthHandlers) Login(c *gin.Context) {
	returnTo := c.Query("returnTo")
	if returnTo == "" {
		returnTo = "/admin"
	}

	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := db.GetUserByEmail(h.db, email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"error":    "Invalid credentials",
			"returnTo": returnTo,
		})
		return
	}

	if err := db.UpdateUserSessionToken(h.db, user); err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.SetCookie("session", user.SessionToken, 3600*24, "/", "", false, true)

	// Redirect to original destination
	c.Redirect(http.StatusFound, returnTo)
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	// Clear session cookie
	c.SetCookie("session", "", -1, "/", "", false, true)

	// Save theme preference before logout
	theme, _ := c.Cookie("admin_theme")

	c.Redirect(http.StatusFound, "/")

	// Restore theme after redirect is set
	if theme != "" {
		c.SetCookie("admin_theme", theme, 3600*24*365, "/", "", false, false)
	}
}
