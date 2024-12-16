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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupPublicRouter creates a test router with embedded templates
func setupPublicRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.SetFuncMap(utils.GetTemplateFuncs())

	// Create minimal templates for testing
	templates := template.Must(template.New("post.tmpl").Parse(`
		<article>
			<h1>{{ .post.Title }}</h1>
			<div>{{ .post.Content }}</div>
			<div>By: {{if .post.Author}}{{.post.Author.FirstName}} {{.post.Author.LastName}}{{else}}<em>Deleted User</em>{{end}}</div>
		</article>
	`))

	// Add error templates
	template.Must(templates.New("404.tmpl").Parse(`<h1>Not Found</h1>`))
	template.Must(templates.New("500.tmpl").Parse(`<h1>Internal Server Error</h1>`))
	template.Must(templates.New("posts.tmpl").Parse(`
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

	router.SetHTMLTemplate(templates)
	return router
}

func TestPostHandlers_GetPostBySlug(t *testing.T) {
	database := db.SetupTestDB()
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	handlers := NewPublicHandlers(repository.NewRepositories(database), cfg)

	// Setup router with test templates
	router := setupPublicRouter()

	// Register handler
	router.GET("/posts/:slug", handlers.GetPostBySlug)

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
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/posts/"+tt.slug, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantStatus == http.StatusOK {
				body := w.Body.String()
				if tt.wantAuthor {
					assert.Contains(t, body, "Test Author")
					assert.NotContains(t, body, "Deleted User")
				} else {
					assert.Contains(t, body, "Deleted User")
					assert.NotContains(t, body, "Test Author")
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
	router := setupPublicRouter()

	// Register handler
	router.GET("/posts", handlers.ListPosts)

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
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that both posts are present with correct author display
	body := w.Body.String()
	assert.Contains(t, body, "Test Post With Author")
	assert.Contains(t, body, "Test Author")
	assert.Contains(t, body, "Test Post Without Author")
	assert.Contains(t, body, "Deleted User")

	// Check pagination data
	assert.Contains(t, body, "Page 1 of 1")
}
