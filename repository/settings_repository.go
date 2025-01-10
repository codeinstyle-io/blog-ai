package repository

import (
	"errors"

	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/system"

	"gorm.io/gorm"
)

type settingsRepository struct {
	db *gorm.DB
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(db *gorm.DB) models.SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) Get() (*models.Settings, error) {
	var settings models.Settings
	err := r.db.First(&settings).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		settings = models.Settings{
			Title:        system.DefaultTitle,
			Subtitle:     system.DefaultSubtitle,
			ChromaStyle:  system.DefaultChromaStyle,
			Theme:        system.DefaultTheme,
			PostsPerPage: system.DefaultPostsPerPage,
		}
		if err := r.Create(settings); err != nil {
			return nil, err
		}
	}

	return &settings, nil
}

func (r *settingsRepository) Create(settings models.Settings) error {
	return r.db.Create(&settings).Error
}

func (r *settingsRepository) Update(settings *models.Settings) error {
	return r.db.Save(settings).Error
}
