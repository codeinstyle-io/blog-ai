package server

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/handlers"
	"codeinstyle.io/captain/middleware"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server represents the HTTP server and its dependencies
type Server struct {
	router     *gin.Engine
	db         *gorm.DB
	repos      *repository.Repositories
	config     *config.Config
	embeddedFS embed.FS
}

// New creates a new server instance
func New(database *gorm.DB, cfg *config.Config, embeddedFS embed.FS) *Server {
	return &Server{
		router:     gin.Default(),
		db:         database,
		repos:      repository.NewRepositories(database),
		config:     cfg,
		embeddedFS: embeddedFS,
	}
}

// loadTemplates loads a glob of templates
func (s *Server) loadTemplates(pattern string) (*template.Template, error) {
	tmpl := template.New("")
	tmpl.Funcs(utils.GetTemplateFuncs())

	// First, get the base directory from the pattern
	baseDir := filepath.Dir(pattern)

	// Walk through all files in the embedded filesystem
	var templates []string
	err := fs.WalkDir(s.embeddedFS, baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Check if file matches *.tmpl pattern
		if filepath.Ext(path) == ".tmpl" {
			templates = append(templates, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking templates directory: %v", err)
	}

	// Load and parse each template
	for _, t := range templates {
		content, err := fs.ReadFile(s.embeddedFS, t)
		if err != nil {
			return nil, fmt.Errorf("error reading template %s: %v", t, err)
		}

		name := filepath.Base(t)
		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %v", t, err)
		}
	}

	return tmpl, nil
}

// setupRouter configures all routes and middleware
func (s *Server) setupRouter(theme *Theme) error {
	// Load embedded public templates as base
	baseTemplates, err := s.loadTemplates("embedded/public/templates")
	if err != nil {
		return fmt.Errorf("error loading public templates: %v", err)
	}

	// Load and merge admin templates
	adminTemplates, err := s.loadTemplates("embedded/admin/templates")
	if err != nil {
		return fmt.Errorf("error loading admin templates: %v", err)
	}

	// Merge theme templates into base templates
	mergedTemplates := append(theme.Templates.Templates(), adminTemplates.Templates()...)

	for _, tmpl := range mergedTemplates {
		if tmpl.Name() == "" {
			continue
		}

		_, err = baseTemplates.AddParseTree(tmpl.Name(), tmpl.Tree)
		if err != nil {
			return fmt.Errorf("error adding template %s: %v", tmpl.Name(), err)
		}
	}

	// Set the combined templates
	s.router.SetHTMLTemplate(baseTemplates)

	// Serve embedded admin static files
	adminStatic, err := fs.Sub(s.embeddedFS, "embedded/admin/static")
	if err != nil {
		return fmt.Errorf("error setting up admin static files: %v", err)
	}
	s.router.StaticFS("/admin/static", http.FS(adminStatic))

	// Serve theme static files
	s.router.StaticFS("/static", http.FS(theme.StaticFS))

	// Register routes
	handlers.RegisterPublicRoutes(s.router, s.repos, s.config)
	handlers.RegisterAdminRoutes(s.router, s.repos, s.config)
	handlers.RegisterAuthRoutes(s.router, s.repos, s.config)

	// Add middleware to load menu items
	s.router.Use(middleware.LoadMenuItems(s.repos))

	return nil
}

// Run starts the HTTP server
func (s *Server) Run() error {
	// Load theme based on config
	var theme *Theme
	var err error

	if s.config.Site.Theme == "" {
		// Use embedded default theme
		theme, err = s.loadEmbeddedTheme()
	} else {
		// Try to load external theme
		themePath := filepath.Join("themes", s.config.Site.Theme)
		theme, err = s.loadExternalTheme(themePath)
	}
	if err != nil {
		return fmt.Errorf("error loading theme: %v", err)
	}

	// Setup router with theme
	if err := s.setupRouter(theme); err != nil {
		return fmt.Errorf("error setting up router: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	fmt.Printf("Server running on http://%s\n", addr)
	return s.router.Run(addr)
}

// InitDevDB initializes the development database with test data
func (s *Server) InitDevDB() error {
	return db.InsertTestData(s.db)
}
