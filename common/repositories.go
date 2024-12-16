package common

import "codeinstyle.io/captain/models"

// Repositories represents all available repositories
type Repositories struct {
	Post     models.PostRepository
	Tag      models.TagRepository
	User     models.UserRepository
	Page     models.PageRepository
	MenuItem models.MenuItemRepository
	Settings models.SettingsRepository
	Media    models.MediaRepository
}
