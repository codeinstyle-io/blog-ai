package handlers

import (
	"net/http"
	"strconv"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"github.com/gin-gonic/gin"
)

// AdminHandlers handles all admin routes
type AdminHandlers struct {
	*BaseHandler
}

// NewAdminHandlers creates a new admin handlers instance
func NewAdminHandlers(repos *repository.Repositories, cfg *config.Config) *AdminHandlers {
	return &AdminHandlers{
		BaseHandler: NewBaseHandler(repos, cfg),
	}
}

func (h *AdminHandlers) Index(c *gin.Context) {
	posts, err := h.repos.Posts.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}
	postCount := int64(len(posts))

	tags, err := h.repos.Tags.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}
	tagCount := int64(len(tags))

	users, err := h.repos.Users.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}
	userCount := int64(len(users))

	// Get 5 most recent posts
	recentPosts, err := h.repos.Posts.FindRecent(5)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
		return
	}

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

func (h *AdminHandlers) ShowSettings(c *gin.Context) {
	settings, err := h.repos.Settings.Get()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(c, gin.H{}))
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

func (h *AdminHandlers) UpdateSettings(c *gin.Context) {
	form, _ := h.repos.Settings.Get()
	var errors []string

	// Get form values
	form.Title = c.PostForm("title")
	form.Subtitle = c.PostForm("subtitle")
	form.Timezone = c.PostForm("timezone")
	form.ChromaStyle = c.PostForm("chroma_style")
	form.Theme = c.PostForm("theme")
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

	// Validate theme
	if form.Theme != "" && form.Theme != "light" && form.Theme != "dark" {
		errors = append(errors, "Invalid theme selected")
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
			"theme":        form.Theme,
			"postsPerPage": form.PostsPerPage,
			"errors":       errors,
		}
		c.HTML(http.StatusBadRequest, "admin_settings.tmpl", h.addCommonData(c, data))
		return
	}

	// Set defaults for optional fields if not provided
	if form.Timezone == "" {
		form.Timezone = system.DefaultTimezone
	}
	if form.ChromaStyle == "" {
		form.ChromaStyle = system.DefaultChromaStyle
	}
	if form.Theme == "" {
		form.Theme = system.DefaultTheme
	}
	if form.PostsPerPage == 0 {
		form.PostsPerPage = system.DefaultPostsPerPage
	}

	if err := h.repos.Settings.Update(form); err != nil {
		errors = append(errors, "Failed to update settings")
		data := gin.H{
			"settings":     form,
			"timezones":    h.config.GetTimezones(),
			"chromaStyles": h.config.GetChromaStyles(),
			"theme":        form.Theme,
			"postsPerPage": form.PostsPerPage,
			"errors":       errors,
		}
		c.HTML(http.StatusInternalServerError, "admin_settings.tmpl", h.addCommonData(c, data))
		return
	}

	c.Redirect(http.StatusFound, "/admin/settings")
}

func (h *AdminHandlers) addCommonData(c *gin.Context, data gin.H) gin.H {
	settings, _ := h.repos.Settings.Get()

	data["settings"] = settings
	return data
}
