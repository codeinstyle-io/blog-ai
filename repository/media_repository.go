package repository

import (
	"github.com/captain-corp/captain/models"

	"gorm.io/gorm"
)

type mediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository creates a new media repository
func NewMediaRepository(db *gorm.DB) models.MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(media *models.Media) error {
	return r.db.Create(media).Error
}

func (r *mediaRepository) Update(media *models.Media) error {
	return r.db.Save(media).Error
}

func (r *mediaRepository) Delete(media *models.Media) error {
	return r.db.Delete(&models.Media{}, media).Error
}

func (r *mediaRepository) FindByID(id uint) (*models.Media, error) {
	var media models.Media
	err := r.db.First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) FindByPath(path string) (*models.Media, error) {
	var media models.Media
	err := r.db.Where("path = ?", path).First(&media).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) FindAll() ([]*models.Media, error) {
	var media []*models.Media
	err := r.db.Order("created_at desc").Find(&media).Error
	return media, err
}
