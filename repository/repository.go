package repository

import (
	"github.com/captain-corp/captain/models"

	"gorm.io/gorm"
)

// Repositories holds all repository implementations
type Repositories struct {
	Posts     models.PostRepository
	Tags      models.TagRepository
	Users     models.UserRepository
	Pages     models.PageRepository
	MenuItems models.MenuItemRepository
	Settings  models.SettingsRepository
	Media     models.MediaRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Posts:     NewPostRepository(db),
		Tags:      NewTagRepository(db),
		Users:     NewUserRepository(db),
		Pages:     NewPageRepository(db),
		MenuItems: NewMenuItemRepository(db),
		Settings:  NewSettingsRepository(db),
		Media:     NewMediaRepository(db),
	}
}
