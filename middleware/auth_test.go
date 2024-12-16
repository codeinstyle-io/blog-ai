package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeinstyle.io/captain/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAuthRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupAuth      func(*gin.Context)
		checkResponse  func(*httptest.ResponseRecorder)
		expectedStatus int
	}{
		{
			name: "Valid session",
			setupAuth: func(c *gin.Context) {
				c.SetCookie("session", "valid-token", 3600, "/", "", false, true)
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "No session cookie",
			setupAuth: func(c *gin.Context) {},
			checkResponse: func(w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusFound, w.Code)
				assert.Equal(t, "/login", w.Header().Get("Location"))
			},
			expectedStatus: http.StatusFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			database := setupTestDB(t)
			repos := handlers.NewRepositories(database)

			router.Use(AuthRequired(repos))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			tt.setupAuth(gin.New().Context(req))

			router.ServeHTTP(w, req)
			tt.checkResponse(w)
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	// Setup test database connection
	// This is just a mock implementation for the test
	return nil
}
