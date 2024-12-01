package handlers

import (
	"fmt"
	"html/template"
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

// setupTestRouter creates a test router with embedded templates
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.SetFuncMap(utils.GetTemplateFuncs())

	// Create a minimal template for testing
	tmpl := template.Must(template.New("admin_tag_posts.tmpl").Parse(`
		<h1>Posts for tag {{ .tag.Name }}</h1>
		<ul>
		{{ range .posts }}
			<li>{{ .Title }}</li>
		{{ end }}
		</ul>
	`))

	router.SetHTMLTemplate(tmpl)
	return router
}

func TestListPostsByTag(t *testing.T) {
	database := db.SetupTestDB()
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	handlers := NewAdminHandlers(database, cfg)

	// Setup router with test templates
	router := setupTestRouter()

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

	// Make request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/tags/%d/posts", tag.ID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Post")
}
