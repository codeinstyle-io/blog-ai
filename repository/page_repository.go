package repository

import (
	"codeinstyle.io/captain/models"
	"gorm.io/gorm"
)

type pageRepository struct {
	db *gorm.DB
}

// NewPageRepository creates a new page repository
func NewPageRepository(db *gorm.DB) models.PageRepository {
	return &pageRepository{db: db}
}

func (r *pageRepository) Create(page *models.Page) error {
	return r.db.Create(page).Error
}

func (r *pageRepository) Update(page *models.Page) error {
	return r.db.Save(page).Error
}

func (r *pageRepository) Delete(page *models.Page) error {
	return r.db.Delete(&models.Page{}, page).Error
}

func (r *pageRepository) FindByID(id uint) (*models.Page, error) {
	var page models.Page
	err := r.db.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *pageRepository) FindAll() ([]*models.Page, error) {
	var pages []*models.Page
	err := r.db.Find(&pages).Error
	return pages, err
}

func (r *pageRepository) FindBySlug(slug string) (*models.Page, error) {
	var page models.Page
	err := r.db.Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}
