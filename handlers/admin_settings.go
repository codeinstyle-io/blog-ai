package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/captain-corp/captain/config"
	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
)

// ShowSettings handles the GET /admin/settings route
func (h *AdminHandlers) ShowSettings(c *fiber.Ctx) error {
	settings, err := h.repos.Settings.Get()
	if err != nil {
		flash.Error(c, "Failed to load settings")
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	var logo *models.Media
	if settings.LogoID != nil {
		logo, _ = h.repos.Media.FindByID(*settings.LogoID)
	}

	data := fiber.Map{
		"title":        "Site Settings",
		"settings":     settings,
		"logo":         logo,
		"timezones":    config.GetTimezones(),
		"chromaStyles": config.GetChromaStyles(),
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
	logoID := c.FormValue("logo_id")
	useFavicon := c.FormValue("use_favicon") == "on"

	// Validate required fields
	if form.Title == "" {
		errors = append(errors, "Title is required")
	}
	if form.Subtitle == "" {
		errors = append(errors, "Subtitle is required")
	}

	// Parse posts per page
	if pp, err := strconv.Atoi(postsPerPage); err == nil {
		form.PostsPerPage = pp
	} else {
		errors = append(errors, "Posts per page must be a number")
	}

	// Handle logo
	if logoID != "" {
		if id, err := strconv.ParseUint(logoID, 10, 32); err == nil {
			uid := uint(id)
			form.LogoID = &uid
		}
	} else {
		form.LogoID = nil
	}
	form.UseFavicon = useFavicon

	if len(errors) > 0 {
		for _, err := range errors {
			flash.Error(c, err)
		}
		return c.Redirect("/admin/settings")
	}

	if err := h.repos.Settings.Update(form); err != nil {
		flash.Error(c, "Failed to save settings")
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate favicons if enabled and logo is set
	if form.UseFavicon && form.LogoID != nil {
		logo, err := h.repos.Media.FindByID(*form.LogoID)
		if err != nil {
			return fmt.Errorf("failed to get logo: %w", err)
		}

		if err := GenerateFavicons(logo, h.storage); err != nil {
			flash.Error(c, "Failed to generate favicons")
			fmt.Printf("Failed to generate favicons: %v\n", err)
			return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
				"error": err.Error(),
			})
		}
	}

	flash.Success(c, "Settings updated successfully")
	return c.Redirect("/admin/settings")
}
