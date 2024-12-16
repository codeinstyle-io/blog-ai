package repository

import (
	"codeinstyle.io/captain/models"
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
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *settingsRepository) Update(settings *models.Settings) error {
	// If no settings exist, create new ones
	var count int64
	r.db.Model(&models.Settings{}).Count(&count)
	if count == 0 {
		return r.db.Create(settings).Error
	}
	// Otherwise update existing settings
	return r.db.Model(&models.Settings{}).Updates(settings).Error
}
