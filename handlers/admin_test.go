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
	cfg := config.NewDefaultConfig()
	handlers := NewAdminHandlers(database, cfg)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.SetFuncMap(utils.GetTemplateFuncs())
	router.LoadHTMLGlob("../templates/**/*.tmpl")

	// Create test data
	tag := db.Tag{Name: "test-tag"}
	database.Create(&tag)

	post := db.Post{
		Title:       "Test Post",
		Slug:        "test-post",
		Content:     "Test content",
		PublishedAt: time.Now(),
		Visible:     true,
		Tags:        []db.Tag{tag},
	}
	database.Create(&post)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Test handler
	handlers.ListPostsByTag(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Post")
	assert.Contains(t, w.Body.String(), "test-tag")
}
