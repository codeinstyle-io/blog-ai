package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequired(t *testing.T) {
	database := db.SetupTestDB()
	router := gin.Default()

	// Create test user with session
	user := &db.User{
		Email:        "test@example.com",
		SessionToken: "valid-token",
	}
	database.Create(user)

	tests := []struct {
		name         string
		token        string
		wantStatus   int
		wantRedirect string
	}{
		{
			name:         "No token",
			token:        "",
			wantStatus:   http.StatusFound,
			wantRedirect: "/login",
		},
		{
			name:       "Valid token",
			token:      "valid-token",
			wantStatus: http.StatusOK,
		},
		{
			name:         "Invalid token",
			token:        "invalid-token",
			wantStatus:   http.StatusFound,
			wantRedirect: "/login",
		},
	}

	// Protected route
	router.GET("/protected", AuthRequired(database), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.token != "" {
				req.AddCookie(&http.Cookie{
					Name:  "session",
					Value: tt.token,
				})
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantRedirect != "" {
				assert.Contains(t, w.Header().Get("Location"), tt.wantRedirect)
			}
		})
	}
}
