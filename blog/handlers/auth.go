package handlers

import (
	"net/http"

	"codeinstyle.io/blog/db"
	"codeinstyle.io/blog/utils"
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
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := db.GetUserByEmail(h.db, email)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if err := db.UpdateUserSessionToken(h.db, user); err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.SetCookie("session", user.SessionToken, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandlers) Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
