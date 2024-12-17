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
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter creates a test router with embedded templates
func setupTestRouter(repositories *repository.Repositories) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: setupTestTemplates(),
	})

	// Create a minimal template for testing
	tmpl := template.Must(template.New("admin_tag_posts").Parse(`
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

	template.Must(tmpl.New("login").Parse(`
		<form method="post" action="/login">
			<input type="email" name="email" />
			<input type="password" name="password" />
			<button type="submit">Login</button>
		</form>
	`))

	// Add the 500 error template
	template.Must(tmpl.New("500").Parse(`
		<h1>Internal Server Error</h1>
		<p>Something went wrong.</p>
	`))

	// Add the 404 error template
	template.Must(tmpl.New("404").Parse(`
		<h1>Not Found</h1>
		<p>The requested resource was not found.</p>
	`))

	// Add the common data template
	template.Must(tmpl.New("common_data").Parse(`
		{{define "common_data"}}
		{{end}}
	`))

	app.Views = tmpl
	return app
}

func setupTestAdmin(t *testing.T) (*fiber.App, *AdminHandlers, *repository.Repositories) {
	cfg := config.NewTestConfig()
	repos := repository.NewTestRepositories()

	app := setupTestRouter(repos)
	handlers := NewAdminHandlers(repos, cfg)

	return app, handlers, repos
}

func TestListPostsByTag(t *testing.T) {
	database := db.SetupTestDB()
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	repos := repository.NewRepositories(database)
	handlers := NewAdminHandlers(repos, cfg)

	// Setup router with test templates
	app := setupTestRouter(repos)

	// Register the handler
	app.Get("/admin/tags/:id/posts", handlers.ListPostsByTag)

	// Create test data
	tag := models.Tag{Name: "test-tag"}
	database.Create(&tag)

	author := models.User{
		FirstName: "Test",
		LastName:  "Author",
		Email:     "test@example.com",
	}
	database.Create(&author)

	// Create settings with timezone
	settings := models.Settings{
		Timezone: "UTC",
	}
	database.Create(&settings)

	// Create post with author
	postWithAuthor := models.Post{
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
	postWithoutAuthor := models.Post{
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
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admin/tags/%d/posts", tag.ID), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := resp.Body.String()
	assert.Contains(t, body, "Test Post With Author")
	assert.Contains(t, body, "Test Author")
	assert.Contains(t, body, "Test Post Without Author")
	assert.Contains(t, body, "Deleted User")
}
