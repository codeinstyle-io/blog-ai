package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"codeinstyle.io/captain/utils"
)

// Theme represents a blog theme with templates and static files
type Theme struct {
	Name      string
	Templates *template.Template
	StaticFS  fs.FS
}

// loadEmbeddedTheme loads the default theme from embedded filesystem
func (s *Server) loadEmbeddedTheme() (*Theme, error) {
	// Load templates
	tmpl, err := s.loadTemplates("embedded/themes/default/templates")
	if err != nil {
		return nil, err
	}

	// Get static files
	staticFS, err := fs.Sub(s.embeddedFS, "embedded/themes/default/static")
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
func (s *Server) loadExternalTheme(themePath string) (*Theme, error) {
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
