package repository

import (
	"codeinstyle.io/captain/models"
	"gorm.io/gorm"
)

type menuItemRepository struct {
	db *gorm.DB
}

// NewMenuItemRepository creates a new menu item repository
func NewMenuItemRepository(db *gorm.DB) models.MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) Create(menuItem *models.MenuItem) error {
	return r.db.Create(menuItem).Error
}

func (r *menuItemRepository) Update(menuItem *models.MenuItem) error {
	return r.db.Save(menuItem).Error
}

func (r *menuItemRepository) Delete(media *models.MenuItem) error {
	return r.db.Delete(&models.MenuItem{}, media).Error
}

func (r *menuItemRepository) FindByID(id uint) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	err := r.db.First(&menuItem, id).Error
	if err != nil {
		return nil, err
	}
	return &menuItem, nil
}

func (r *menuItemRepository) FindAll() ([]*models.MenuItem, error) {
	var menuItems []*models.MenuItem
	err := r.db.Order("position").Find(&menuItems).Error
	return menuItems, err
}

func (r *menuItemRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&models.MenuItem{}).Error
}

func (r *menuItemRepository) UpdatePositions(startPosition int) error {
	return r.db.Model(&models.MenuItem{}).
		Where("position > ?", startPosition).
		UpdateColumn("position", gorm.Expr("position - 1")).
		Error
}

func (r *menuItemRepository) GetNextPosition() int {
	var maxPosition struct {
		MaxPos int
	}
	r.db.Model(&models.MenuItem{}).
		Select("COALESCE(MAX(position), 0) as max_pos").
		Scan(&maxPosition)
	return maxPosition.MaxPos + 1
}
