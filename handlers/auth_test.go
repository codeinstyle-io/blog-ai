package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"codeinstyle.io/captain/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupAuthTest(t *testing.T, database *gorm.DB) (*fiber.App, *AuthHandlers) {
	repositories := repository.NewRepositories(database)
	app := fiber.New()

	// Create minimal templates for testing
	authHandlers := NewAuthHandlers(repositories, &config.Config{})

	// Create test user
	password := "Test1234!"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	testUser := models.User{
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	err = database.Create(&testUser).Error
	assert.NoError(t, err)

	return app, authHandlers
}

func TestAuthHandlers_Login(t *testing.T) {
	database := db.SetupTestDB()
	app, authHandlers := setupAuthTest(t, database)

	// Add routes
	app.Get("/login", authHandlers.Login)
	app.Post("/login", authHandlers.PostLogin)

	tests := []struct {
		name          string
		email         string
		password      string
		next          string
		expectedCode  int
		expectedPath  string
		sessionCookie bool
	}{
		{
			name:          "Valid login",
			email:         "test@example.com",
			password:      "Test1234!",
			expectedCode:  http.StatusFound,
			expectedPath:  "/admin",
			sessionCookie: true,
		},
		{
			name:          "Invalid password",
			email:         "test@example.com",
			password:      "wrong",
			expectedCode:  http.StatusUnauthorized,
			expectedPath:  "",
			sessionCookie: false,
		},
		{
			name:          "Invalid email",
			email:         "wrong@example.com",
			password:      "Test1234!",
			expectedCode:  http.StatusUnauthorized,
			expectedPath:  "",
			sessionCookie: false,
		},
		{
			name:          "Custom return path",
			email:         "test@example.com",
			password:      "Test1234!",
			next:          "/admin/posts",
			expectedCode:  http.StatusFound,
			expectedPath:  "/admin/posts",
			sessionCookie: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			// Add return path if specified
			if tt.next != "" {
				form.Add("next", tt.next)
			}

			// Create request
			req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte(form.Encode())))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Create response recorder
			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			if tt.expectedPath != "" {
				assert.Equal(t, tt.expectedPath, resp.Header.Get("Location"))
			}

			// Check session cookie
			if tt.sessionCookie {
				assert.Contains(t, resp.Header.Get("Set-Cookie"), system.CookieName+"=")
			} else {
				assert.NotContains(t, resp.Header.Get("Set-Cookie"), system.CookieName+"=")
			}
		})
	}
}

func TestAuthHandlers_Logout(t *testing.T) {
	database := db.SetupTestDB()
	app, authHandlers := setupAuthTest(t, database)

	// Add route
	app.Get("/logout", authHandlers.Logout)

	// Test GET /logout
	req := httptest.NewRequest("GET", "/logout", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/login", resp.Header.Get("Location"))

	// Check that session cookie is cleared
	cookies := resp.Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == system.CookieName {
			sessionCookie = cookie
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.True(t, sessionCookie.MaxAge < 0)
}
