package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
			name: "No session cookie",
			setupAuth: func(c *gin.Context) {
				// Don't set any cookie
			},
			checkResponse: func(w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusFound, w.Code)
				assert.Equal(t, "/login", w.Header().Get("Location"))
			},
			expectedStatus: http.StatusFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := testutils.SetupTestDB(t)
			repos := repository.NewRepositories(db)

			router := gin.New()
			router.Use(AuthRequired(repos))
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.setupAuth != nil {
				ctx, _ := gin.CreateTestContext(w)
				ctx.Request = req
				tt.setupAuth(ctx)
				// Copy cookies from the test context to the request
				for _, cookie := range w.Result().Cookies() {
					req.AddCookie(cookie)
				}
			}

			router.ServeHTTP(w, req)
			if tt.checkResponse != nil {
				tt.checkResponse(w)
			}
		})
	}
}
