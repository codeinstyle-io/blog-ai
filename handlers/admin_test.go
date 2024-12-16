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
	"codeinstyle.io/captain/repository"
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
		{{if .error}}
		<div class="error">{{.error}}</div>
		{{end}}
		<h1>Posts for tag {{ .tag.Name }}</h1>
		<ul>
		{{ range .posts }}
			<li>{{ .Title }} - By: {{if .Author}}{{.Author.FirstName}} {{.Author.LastName}}{{else}}<em>Deleted User</em>{{end}}</li>
		{{ end }}
		</ul>
	`))

	// Add the 500 error template
	template.Must(tmpl.New("500.tmpl").Parse(`
		<h1>Internal Server Error</h1>
		<p>Something went wrong.</p>
	`))

	// Add the 404 error template
	template.Must(tmpl.New("404.tmpl").Parse(`
		<h1>Not Found</h1>
		<p>The requested resource was not found.</p>
	`))

	// Add the common data template
	template.Must(tmpl.New("common_data.tmpl").Parse(`
		{{define "common_data"}}
		{{end}}
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
	handlers := NewAdminHandlers(repository.NewRepositories(database), cfg)

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

	// Create settings with timezone
	settings := db.Settings{
		Timezone: "UTC",
	}
	database.Create(&settings)

	// Create post with author
	postWithAuthor := db.Post{
		Title:       "Test Post With Author",
		Slug:        "test-post-with-author",
		Content:     "Test content",
		PublishedAt: time.Now(),
		Visible:     true,
		AuthorID:    author.ID,
	}
	database.Create(&postWithAuthor)

	// Associate post with tag
	if err := database.Model(&postWithAuthor).Association("Tags").Append(&tag); err != nil {
		t.Fatalf("Failed to associate tag with post: %v", err)
	}

	// Create post without author
	postWithoutAuthor := db.Post{
		Title:       "Test Post Without Author",
		Slug:        "test-post-without-author",
		Content:     "Test content",
		PublishedAt: time.Now(),
		Visible:     true,
	}
	database.Create(&postWithoutAuthor)

	// Associate post with tag
	if err := database.Model(&postWithoutAuthor).Association("Tags").Append(&tag); err != nil {
		t.Fatalf("Failed to associate tag with post: %v", err)
	}

	// Make request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", fmt.Sprintf("/admin/tags/%d/posts", tag.ID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Test Post With Author")
	assert.Contains(t, body, "Test Author")
	assert.Contains(t, body, "Test Post Without Author")
	assert.Contains(t, body, "Deleted User")
}
