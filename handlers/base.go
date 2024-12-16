package handlers

import (
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"gorm.io/gorm"
)

// BaseHandler contains common dependencies for all handlers
type BaseHandler struct {
	config       *config.Config
	repos        *repository.Repositories
	settings     *models.Settings
	postRepo     models.PostRepository
	tagRepo      models.TagRepository
	userRepo     models.UserRepository
	pageRepo     models.PageRepository
	menuRepo     models.MenuItemRepository
	settingsRepo models.SettingsRepository
	mediaRepo    models.MediaRepository
}

// NewBaseHandler creates a new base handler with common dependencies
func NewBaseHandler(repos *repository.Repositories, cfg *config.Config) *BaseHandler {
	settings, _ := repos.Settings.Get()
	return &BaseHandler{
		config:       cfg,
		repos:        repos,
		settings:     settings,
		postRepo:     repos.Posts,
		tagRepo:      repos.Tags,
		userRepo:     repos.Users,
		pageRepo:     repos.Pages,
		menuRepo:     repos.MenuItems,
		settingsRepo: repos.Settings,
		mediaRepo:    repos.Media,
	}
}

// NewRepositories creates a new instance of Repositories
func NewRepositories(db interface{}) *repository.Repositories {
	gormDB := db.(*gorm.DB)
	return &repository.Repositories{
		Posts:     repository.NewPostRepository(gormDB),
		Tags:      repository.NewTagRepository(gormDB),
		Users:     repository.NewUserRepository(gormDB),
		Pages:     repository.NewPageRepository(gormDB),
		MenuItems: repository.NewMenuItemRepository(gormDB),
		Settings:  repository.NewSettingsRepository(gormDB),
		Media:     repository.NewMediaRepository(gormDB),
		Sessions:  repository.NewSessionRepository(gormDB),
	}
}
