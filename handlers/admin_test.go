package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListPostsByTag(t *testing.T) {
	database := db.SetupTestDB()
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	handlers := NewAdminHandlers(database, cfg)

	// Setup Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.SetFuncMap(utils.GetTemplateFuncs())
	router.LoadHTMLGlob("../templates/**/*.tmpl")

	// Register the handler
	router.GET("/admin/tags/:id/posts", handlers.ListPostsByTag)

	// Create test data
	tag := db.Tag{Name: "test-tag"}
	database.Create(&tag)

	author := db.User{
		FirstName: "Test",
		LastName:  "Author",
		Email:     "test@example.com",
	}
	database.Create(&author)

	post := db.Post{
		Title:       "Test Post",
		Slug:        "test-post",
		Content:     "Test content",
		PublishedAt: time.Now(),
		Visible:     true,
		Tags:        []db.Tag{tag},
		AuthorID:    author.ID,
	}
	database.Create(&post)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/tags/1/posts", nil)

	// Serve the request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Post")
	assert.Contains(t, w.Body.String(), "test-tag")
}
