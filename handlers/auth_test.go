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
	gin.SetMode(gin.TestMode)
	router := gin.New() // Don't use Default() to avoid extra middleware

	// Add template functions
	router.SetFuncMap(utils.GetTemplateFuncs())
	router.LoadHTMLGlob("../templates/**/*.tmpl")

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

	router.POST("/login", authHandlers.Login) // Move this outside the test loop

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			// Fix URL construction for return path
			targetUrl := "/login"
			if tt.returnTo != "" {
				targetUrl += "?returnTo=" + tt.returnTo
			}

			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			req := httptest.NewRequest("POST", targetUrl, strings.NewReader(form.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add template functions
	router.SetFuncMap(utils.GetTemplateFuncs())

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

	router.GET("/logout", authHandlers.Logout)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/logout", nil)

			// Set theme cookie if specified
			if tt.theme != "" {
				req.AddCookie(&http.Cookie{
					Name:  "admin_theme",
					Value: tt.theme,
				})
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusFound, w.Code)
			assert.Equal(t, "/", w.Header().Get("Location"))

			var sessionCookie *http.Cookie
			for _, cookie := range w.Result().Cookies() {
				if cookie.Name == "session" {
					sessionCookie = cookie
					break
				}
			}

			// Only verify session cookie was cleared
			assert.NotNil(t, sessionCookie)
			assert.True(t, sessionCookie.MaxAge < 0)
		})
	}
}
