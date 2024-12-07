package handlers

import (
	"net/http"
	"strconv"
	"time"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandlers struct {
	db     *gorm.DB
	config *config.Config
}

func NewAdminHandlers(database *gorm.DB, config *config.Config) *AdminHandlers {
	return &AdminHandlers{
		db:     database,
		config: config,
	}
}

func (h *AdminHandlers) Index(c *gin.Context) {
	var postCount, tagCount, userCount int64
	var recentPosts []db.Post

	// Get counts
	h.db.Model(&db.Post{}).Count(&postCount)
	h.db.Model(&db.Tag{}).Count(&tagCount)
	h.db.Model(&db.User{}).Count(&userCount)

	// Get 5 most recent posts
	h.db.Order("published_at desc").Limit(5).Find(&recentPosts)

	data := gin.H{
		"title":       "Dashboard",
		"postCount":   postCount,
		"tagCount":    tagCount,
		"userCount":   userCount,
		"recentPosts": recentPosts,
	}

	data = h.addCommonData(c, data)

	c.HTML(http.StatusOK, "admin_index.tmpl", data)
}

// handlers/admin.go
func (h *AdminHandlers) SavePreferences(c *gin.Context) {
	var prefs struct {
		Theme string `json:"theme"`
	}

	if err := c.BindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preferences"})
		return
	}

	// Save theme preference in cookie
	c.SetCookie("admin_theme", prefs.Theme, 3600*24*365, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *AdminHandlers) ShowSettings(c *gin.Context) {
	settings, err := db.GetSettings(h.db)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	data := gin.H{
		"title":        "Site Settings",
		"settings":     settings,
		"timezones":    h.config.GetTimezones(),
		"chromaStyles": h.config.GetChromaStyles(),
	}

	data = h.addCommonData(c, data)
	c.HTML(http.StatusOK, "admin_settings.tmpl", data)
}

const (
	DefaultTimezone     = "UTC"
	DefaultChromaStyle  = "paraiso-dark"
	DefaultPostPerPage  = 10
)

func (h *AdminHandlers) UpdateSettings(c *gin.Context) {
	var form db.Settings
	var errors []string

	// Get form values
	form.Title = c.PostForm("title")
	form.Subtitle = c.PostForm("subtitle")
	form.Timezone = c.PostForm("timezone")
	form.ChromaStyle = c.PostForm("chroma_style")
	postsPerPage := c.PostForm("posts_per_page")

	// Validate required fields
	if form.Title == "" {
		errors = append(errors, "Title is required")
	}
	if form.Subtitle == "" {
		errors = append(errors, "Subtitle is required")
	}

	// Validate timezone
	if form.Timezone != "" {
		valid := false
		for _, tz := range h.config.GetTimezones() {
			if tz == form.Timezone {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, "Invalid timezone selected")
		}
	}

	// Validate chroma style
	if form.ChromaStyle != "" {
		valid := false
		for _, style := range h.config.GetChromaStyles() {
			if style == form.ChromaStyle {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, "Invalid syntax highlighting theme selected")
		}
	}

	// Parse and validate posts per page
	if postsPerPage != "" {
		if pp, err := strconv.Atoi(postsPerPage); err != nil {
			errors = append(errors, "Posts per page must be a number")
		} else if pp < 1 || pp > 50 {
			errors = append(errors, "Posts per page must be between 1 and 50")
		} else {
			form.PostsPerPage = pp
		}
	}

	if len(errors) > 0 {
		data := gin.H{
			"settings":     form,
			"timezones":    h.config.GetTimezones(),
			"chromaStyles": h.config.GetChromaStyles(),
			"errors":       errors,
		}
		c.HTML(http.StatusBadRequest, "admin_settings.tmpl", data)
		return
	}

	// Set defaults for optional fields if not provided
	if form.Timezone == "" {
		form.Timezone = DefaultTimezone
	}
	if form.ChromaStyle == "" {
		form.ChromaStyle = DefaultChromaStyle
	}
	if form.PostsPerPage == 0 {
		form.PostsPerPage = DefaultPostPerPage
	}

	form.LastUpdatedAt = time.Now()

	if err := db.UpdateSettings(h.db, &form); err != nil {
		errors = append(errors, "Failed to update settings")
		data := gin.H{
			"settings":     form,
			"timezones":    h.config.GetTimezones(),
			"chromaStyles": h.config.GetChromaStyles(),
			"errors":       errors,
		}
		c.HTML(http.StatusInternalServerError, "admin_settings.tmpl", data)
		return
	}

	c.Redirect(http.StatusFound, "/admin/settings")
}
