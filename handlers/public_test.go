package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// setupPublicRouter creates a test router with embedded templates
func setupPublicRouter() *fiber.App {
	app := fiber.New()
	app.Settings.TemplateDir = "./templates"
	app.Settings.TemplateExtension = ""

	// Create minimal templates for testing
	templates := template.Must(template.New("post").Parse(`
		<article>
			<h1>{{ .post.Title }}</h1>
			<div>{{ .post.Content }}</div>
			<div>By: {{if .post.Author}}{{.post.Author.FirstName}} {{.post.Author.LastName}}{{else}}<em>Deleted User</em>{{end}}</div>
		</article>
	`))

	// Add error templates
	template.Must(templates.New("404").Parse(`<h1>Not Found</h1>`))
	template.Must(templates.New("500").Parse(`<h1>Internal Server Error</h1>`))
	template.Must(templates.New("posts").Parse(`
		<div>
		{{ range .posts }}
			<article>
				<h2>{{ .Title }}</h2>
				<div>{{ .Excerpt }}</div>
				<div>By: {{if .Author}}{{.Author.FirstName}} {{.Author.LastName}}{{else}}<em>Deleted User</em>{{end}}</div>
			</article>
		{{ end }}
		<div class="pagination">
			Page {{ .currentPage }} of {{ .totalPages }}
		</div>
		</div>
	`))

	app.Settings.TemplateFileSystem = utils.GetTemplateFuncs()
	return app
}

func TestPostHandlers_GetPostBySlug(t *testing.T) {
	database := db.SetupTestDB()
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	handlers := NewPublicHandlers(repository.NewRepositories(database), cfg)

	// Setup router with test templates
	app := setupPublicRouter()

	// Register handler
	app.Get("/posts/:slug", handlers.GetPostBySlug)

	// Create test author
	author := &models.User{
		FirstName: "Test",
		LastName:  "Author",
		Email:     "test@example.com",
	}
	database.Create(author)

	// Create test post with author
	postWithAuthor := &models.Post{
		Title:       "Test Post With Author",
		Slug:        "test-post-with-author",
		Content:     "Test Content",
		PublishedAt: time.Now(),
		Visible:     true,
		AuthorID:    author.ID,
	}
	database.Create(postWithAuthor)

	// Create test post without author
	postWithoutAuthor := &models.Post{
		Title:       "Test Post Without Author",
		Slug:        "test-post-without-author",
		Content:     "Test Content",
		PublishedAt: time.Now(),
		Visible:     true,
	}
	database.Create(postWithoutAuthor)

	tests := []struct {
		name       string
		slug       string
		wantStatus int
		wantAuthor bool
	}{
		{
			name:       "Post with author",
			slug:       "test-post-with-author",
			wantStatus: http.StatusOK,
			wantAuthor: true,
		},
		{
			name:       "Post without author",
			slug:       "test-post-without-author",
			wantStatus: http.StatusOK,
			wantAuthor: false,
		},
		{
			name:       "Non-existent post",
			slug:       "missing",
			wantStatus: http.StatusNotFound,
			wantAuthor: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/posts/"+tt.slug, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			if tt.wantStatus == http.StatusOK {
				body := resp.Body()
				if tt.wantAuthor {
					assert.Contains(t, body.String(), "Test Author")
					assert.NotContains(t, body.String(), "Deleted User")
				} else {
					assert.Contains(t, body.String(), "Deleted User")
					assert.NotContains(t, body.String(), "Test Author")
				}
			}
		})
	}
}

func TestPostHandlers_ListPosts(t *testing.T) {
	database := db.SetupTestDB()
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create default settings first
	settings := &models.Settings{
		PostsPerPage: 2, // Set a small number to test pagination
		Title:        "Test Blog",
		Subtitle:     "Test Subtitle",
		Theme:        "default",
		ChromaStyle:  "monokai",
		Timezone:     "UTC",
	}
	if err := database.Create(settings).Error; err != nil {
		t.Fatalf("Failed to create settings: %v", err)
	}

	handlers := NewPublicHandlers(repository.NewRepositories(database), cfg)

	// Setup router with test templates
	app := setupPublicRouter()

	// Register handler
	app.Get("/posts", handlers.ListPosts)

	// Create test author
	author := &models.User{
		FirstName: "Test",
		LastName:  "Author",
		Email:     "test@example.com",
	}
	database.Create(author)

	// Create test posts
	now := time.Now()
	posts := []models.Post{
		{
			Title:       "Test Post With Author",
			Slug:        "test-post-1",
			Content:     "# Test Content 1",
			PublishedAt: now,
			Visible:     true,
			AuthorID:    author.ID,
		},
		{
			Title:       "Test Post Without Author",
			Slug:        "test-post-2",
			Content:     "# Test Content 2",
			PublishedAt: now,
			Visible:     true,
		},
	}

	// Create posts in database and check for errors
	for _, post := range posts {
		if err := database.Create(&post).Error; err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}
	}

	// Test listing posts
	req := httptest.NewRequest("GET", "/posts", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that both posts are present with correct author display
	body := resp.Body()
	assert.Contains(t, body.String(), "Test Post With Author")
	assert.Contains(t, body.String(), "Test Author")
	assert.Contains(t, body.String(), "Test Post Without Author")
	assert.Contains(t, body.String(), "Deleted User")

	// Check pagination data
	assert.Contains(t, body.String(), "Page 1 of 1")
}
