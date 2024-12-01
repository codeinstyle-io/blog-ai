package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"codeinstyle.io/captain/cli"
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/handlers"
	"codeinstyle.io/captain/middleware"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

//go:embed embedded/admin/static/css/*
//go:embed embedded/admin/static/js/*
//go:embed embedded/admin/templates/*
//go:embed embedded/public/templates/errors/*
//go:embed embedded/public/templates/*
//go:embed embedded/themes/default/static/css/*
//go:embed embedded/themes/default/static/js/*
//go:embed embedded/themes/default/templates/*
var embeddedFS embed.FS

// Theme represents a Captain theme with its templates and static files
type Theme struct {
	Name      string
	Templates *template.Template
	StaticFS  fs.FS
}

// loadTemplates loads a glob of templates
func loadTemplates(pattern string) (*template.Template, error) {
	tmpl := template.New("")

	tmpl.Funcs(utils.GetTemplateFuncs())

	// First, get the base directory from the pattern
	baseDir := filepath.Dir(pattern)

	// Walk through all files in the embedded filesystem
	var templates []string
	err := fs.WalkDir(embeddedFS, baseDir, func(path string, d fs.DirEntry, err error) error {
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
		content, err := fs.ReadFile(embeddedFS, t)
		if err != nil {
			return nil, fmt.Errorf("error reading template %s: %v", t, err)
		}

		// Use the full path as the template name to maintain uniqueness
		name := filepath.Base(t)

		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %v", t, err)
		}
	}

	return tmpl, nil
}

// loadEmbeddedTheme loads the default theme from embedded files
func loadEmbeddedTheme() (*Theme, error) {
	// Load templates
	tmpl, err := loadTemplates("embedded/themes/default/templates")
	if err != nil {
		return nil, err
	}

	// Get static files
	staticFS, err := fs.Sub(embeddedFS, "embedded/themes/default/static")
	if err != nil {
		return nil, fmt.Errorf("error setting up embedded static files: %v", err)
	}

	return &Theme{
		Name:      "default",
		Templates: tmpl,
		StaticFS:  staticFS,
	}, nil
}

// loadExternalTheme loads a theme from the filesystem
func loadExternalTheme(themePath string) (*Theme, error) {
	// Check if theme exists
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("theme not found at %s", themePath)
	}

	// Load templates
	tmpl := template.New("")
	tmpl.Funcs(utils.GetTemplateFuncs())
	templatePath := filepath.Join(themePath, "templates")
	templates, err := filepath.Glob(filepath.Join(templatePath, "*.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("error loading theme templates: %v", err)
	}

	for _, t := range templates {
		content, err := os.ReadFile(t)
		if err != nil {
			return nil, fmt.Errorf("error reading template %s: %v", t, err)
		}
		name := filepath.Base(t)
		_, err = tmpl.New(name).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %v", t, err)
		}
	}

	// Create static files FS
	staticFS := os.DirFS(filepath.Join(themePath, "static"))

	return &Theme{
		Name:      filepath.Base(themePath),
		Templates: tmpl,
		StaticFS:  staticFS,
	}, nil
}

func setupRouter(cfg *config.Config, theme *Theme, database *gorm.DB) (*gin.Engine, error) {
	r := gin.Default()

	// Load embedded public templates as base
	baseTemplates, err := loadTemplates("embedded/public/templates")
	if err != nil {
		return nil, fmt.Errorf("error loading public templates: %v", err)
	}

	// Load and merge admin templates
	adminTemplates, err := loadTemplates("embedded/admin/templates")
	if err != nil {
		return nil, fmt.Errorf("error loading admin templates: %v", err)
	}

	// Merge theme templates into base templates
	mergedTemplates := append(theme.Templates.Templates(), adminTemplates.Templates()...)

	for _, tmpl := range mergedTemplates {
		if tmpl.Name() == "" {
			continue
		}

		_, err = baseTemplates.AddParseTree(tmpl.Name(), tmpl.Tree)
		if err != nil {
			return nil, fmt.Errorf("error adding template %s: %v", tmpl.Name(), err)
		}
	}

	// Set the combined templates
	r.SetHTMLTemplate(baseTemplates)

	// Serve embedded admin static files
	adminStatic, err := fs.Sub(embeddedFS, "embedded/admin/static")
	if err != nil {
		return nil, fmt.Errorf("error setting up admin static files: %v", err)
	}
	r.StaticFS("/admin/static", http.FS(adminStatic))

	// Serve theme static files
	r.StaticFS("/static", http.FS(theme.StaticFS))

	// Register routes with config
	handlers.RegisterPublicRoutes(r, database, cfg)
	handlers.RegisterAdminRoutes(r, database, cfg)

	// Add middleware to load menu items
	r.Use(middleware.LoadMenuItems(database))

	return r, nil
}

var (
	initDevDB  bool
	configFile string
)

func main() {
	var rootCmd = &cobra.Command{Use: "captain"}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the server",
		Run:   runServer,
	}

	runCmd.Flags().BoolVarP(&initDevDB, "init-dev-db", "i", false, "Initialize the development database with test data")

	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	var userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run:   cli.CreateUser,
	}

	var userUpdatePasswordCmd = &cobra.Command{
		Use:   "update-password",
		Short: "Update user password",
		Run:   cli.UpdateUserPassword,
	}

	userCmd.AddCommand(userCreateCmd, userUpdatePasswordCmd)
	rootCmd.AddCommand(runCmd, userCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	// Load config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	database := db.InitDB(cfg)

	// Load theme based on config
	var theme *Theme
	if cfg.Site.Theme == "" {
		// Use embedded default theme
		theme, err = loadEmbeddedTheme()
	} else {
		// Try to load external theme
		themePath := filepath.Join("themes", cfg.Site.Theme)
		theme, err = loadExternalTheme(themePath)
	}
	if err != nil {
		log.Fatalf("Error loading theme: %v", err)
	}

	// Setup router
	r, err := setupRouter(cfg, theme, database)
	if err != nil {
		log.Fatalf("Error setting up router: %v", err)
	}

	if initDevDB {
		err := db.InsertTestData(database)
		if err != nil {
			log.Fatalf("failed to insert test data: %v", err)
		}
	}

	fmt.Printf("Server running on http://%s:%d\n", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
