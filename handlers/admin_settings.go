package handlers

import (
	"net/http"
	"strconv"

	"codeinstyle.io/captain/flash"
	"codeinstyle.io/captain/system"
	"github.com/gofiber/fiber/v2"
)

// ShowSettings handles the GET /admin/settings route
func (h *AdminHandlers) ShowSettings(c *fiber.Ctx) error {
	settings, err := h.repos.Settings.Get()
	if err != nil {
		flash.Error(c, "Failed to load settings")
		return c.Status(http.StatusInternalServerError).Render("500", fiber.Map{
			"err": err.Error(),
		})
	}

	data := fiber.Map{
		"title":        "Site Settings",
		"settings":     settings,
		"timezones":    h.config.GetTimezones(),
		"chromaStyles": h.config.GetChromaStyles(),
	}

	return c.Render("admin_settings", data)
}

// UpdateSettings handles the POST /admin/settings route
func (h *AdminHandlers) UpdateSettings(c *fiber.Ctx) error {
	form, _ := h.repos.Settings.Get()
	var errors []string

	// Get form values
	form.Title = c.FormValue("title")
	form.Subtitle = c.FormValue("subtitle")
	form.Timezone = c.FormValue("timezone")
	form.ChromaStyle = c.FormValue("chroma_style")
	form.Theme = c.FormValue("theme")
	postsPerPage := c.FormValue("posts_per_page")

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
		for _, err := range errors {
			flash.Error(c, err)
		}
		data := fiber.Map{
			"title":        "Site Settings",
			"settings":     form,
			"timezones":    h.config.GetTimezones(),
			"chromaStyles": h.config.GetChromaStyles(),
			"theme":        form.Theme,
			"postsPerPage": form.PostsPerPage,
		}
		return c.Status(http.StatusBadRequest).Render("admin_settings", data)
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
		flash.Error(c, "Failed to update settings")
		data := fiber.Map{
			"title":        "Site Settings",
			"settings":     form,
			"timezones":    h.config.GetTimezones(),
			"chromaStyles": h.config.GetChromaStyles(),
			"theme":        form.Theme,
			"postsPerPage": form.PostsPerPage,
		}

		return c.Status(http.StatusInternalServerError).Render("admin_settings", data)
	}

	flash.Success(c, "Settings updated successfully")
	return c.Redirect("/admin/settings")
}
