package handlers

import (
	"captain-corp/captain/config"
	"captain-corp/captain/models"
	"captain-corp/captain/repository"

	"github.com/gofiber/fiber/v2"
)

// BaseHandlers contains common handler dependencies
type BaseHandlers struct {
	repos  *repository.Repositories
	config *config.Config
}

// NewBaseHandlers creates a new BaseHandlers instance
func NewBaseHandlers(repos *repository.Repositories, cfg *config.Config) *BaseHandlers {
	return &BaseHandlers{
		repos:  repos,
		config: cfg,
	}
}

// addCommonData adds common template data
func (h *BaseHandlers) addCommonData(c *fiber.Ctx, data fiber.Map) fiber.Map {
	if data == nil {
		data = fiber.Map{}
	}

	// Get current user if authenticated
	if auth := c.Locals("user"); auth != nil {
		if user, ok := auth.(*models.User); ok {
			data["currentUser"] = user
		}
	}

	return data
}
