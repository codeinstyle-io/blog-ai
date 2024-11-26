package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlers_Login(t *testing.T) {
	database := db.SetupTestDB()
	router := gin.Default()
	authHandlers := NewAuthHandlers(database)

	// Create test user
	password := "Test1234!"
	hash, _ := utils.HashPassword(password)
	user := &db.User{
		Email:    "test@example.com",
		Password: hash,
	}
	database.Create(user)

	tests := []struct {
		name         string
		email        string
		password     string
		returnTo     string
		wantStatus   int
		wantRedirect string
	}{
		{
			name:         "Valid login",
			email:        "test@example.com",
			password:     "Test1234!",
			wantStatus:   http.StatusFound,
			wantRedirect: "/admin",
		},
		{
			name:       "Invalid password",
			email:      "test@example.com",
			password:   "wrong",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid email",
			email:      "wrong@example.com",
			password:   "Test1234!",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:         "Custom return path",
			email:        "test@example.com",
			password:     "Test1234!",
			returnTo:     "/admin/posts",
			wantStatus:   http.StatusFound,
			wantRedirect: "/admin/posts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			req := httptest.NewRequest("POST", "/login"+tt.returnTo, strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			router.POST("/login", authHandlers.Login)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantRedirect != "" {
				assert.Equal(t, tt.wantRedirect, w.Header().Get("Location"))
			}
		})
	}
}

func TestAuthHandlers_Logout(t *testing.T) {
	database := db.SetupTestDB()
	router := gin.Default()
	authHandlers := NewAuthHandlers(database)

	tests := []struct {
		name  string
		theme string
	}{
		{
			name: "Logout without theme",
		},
		{
			name:  "Logout with theme",
			theme: "dark",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/logout", nil)

			if tt.theme != "" {
				req.AddCookie(&http.Cookie{
					Name:  "admin_theme",
					Value: tt.theme,
				})
			}

			router.GET("/logout", authHandlers.Logout)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusFound, w.Code)
			assert.Equal(t, "/", w.Header().Get("Location"))

			var sessionCookie, themeCookie *http.Cookie
			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == "session" {
					sessionCookie = cookie
				}
				if cookie.Name == "admin_theme" {
					themeCookie = cookie
				}
			}

			assert.NotNil(t, sessionCookie)
			assert.Less(t, sessionCookie.MaxAge, 0)

			if tt.theme != "" {
				assert.NotNil(t, themeCookie)
				assert.Equal(t, tt.theme, themeCookie.Value)
			}
		})
	}
}
