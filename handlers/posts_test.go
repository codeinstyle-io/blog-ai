package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostHandlers_GetPostBySlug(t *testing.T) {
	database := db.SetupTestDB()
	handlers := NewPostHandlers(database)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.SetFuncMap(utils.GetTemplateFuncs())
	router.LoadHTMLGlob("../templates/**/*.tmpl")

	// Create test post
	post := &db.Post{
		Title:       "Test Post",
		Slug:        "test-post",
		Content:     "Test Content",
		PublishedAt: time.Now(),
		Visible:     true,
	}
	database.Create(post)

	tests := []struct {
		name       string
		slug       string
		wantStatus int
	}{
		{
			name:       "Existing post",
			slug:       "test-post",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Non-existent post",
			slug:       "missing",
			wantStatus: http.StatusNotFound,
		},
	}

	router.GET("/posts/:slug", handlers.GetPostBySlug)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/posts/"+tt.slug, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestPostHandlers_ListPosts(t *testing.T) {
	database := db.SetupTestDB()
	handlers := NewPostHandlers(database)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.SetFuncMap(utils.GetTemplateFuncs())
	router.LoadHTMLGlob("../templates/**/*.tmpl")

	// Create test posts
	posts := []db.Post{
		{Title: "Post 1", Slug: "post-1", Visible: true, PublishedAt: time.Now()},
		{Title: "Post 2", Slug: "post-2", Visible: true, PublishedAt: time.Now()},
		{Title: "Hidden", Slug: "hidden", Visible: false, PublishedAt: time.Now()},
	}
	for _, p := range posts {
		database.Create(&p)
	}

	tests := []struct {
		name          string
		page          string
		wantStatus    int
		wantPostCount int
	}{
		{
			name:          "First page",
			page:          "1",
			wantStatus:    http.StatusOK,
			wantPostCount: 2, // Only visible posts
		},
		{
			name:          "Empty page",
			page:          "999",
			wantStatus:    http.StatusOK,
			wantPostCount: 0,
		},
	}

	router.GET("/", handlers.ListPosts)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/?page="+tt.page, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
