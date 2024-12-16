package repository

import (
	"errors"

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

func (r *menuItemRepository) CreateAll(menuItems []models.MenuItem) error {
	return r.db.Create(&menuItems).Error
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

func (r *menuItemRepository) Move(id uint, direction string) error {

	// Start transaction
	tx := r.db.Begin()

	var currentItem models.MenuItem
	if err := tx.First(&currentItem, id).Error; err != nil {
		tx.Rollback()

		return errors.New("Menu item not found")
	}

	// Find adjacent item
	var adjacentItem models.MenuItem
	if direction == "up" {
		if err := tx.Where("position < ?", currentItem.Position).Order("position DESC").First(&adjacentItem).Error; err != nil {
			tx.Rollback()
			return errors.New("Item already at top")
		}
	} else {
		if err := tx.Where("position > ?", currentItem.Position).Order("position ASC").First(&adjacentItem).Error; err != nil {
			tx.Rollback()
			return errors.New("Item already at bottom")
		}
	}

	// Swap positions
	currentPos := currentItem.Position
	adjacentPos := adjacentItem.Position

	if err := tx.Model(&currentItem).Update("position", adjacentPos).Error; err != nil {
		tx.Rollback()
		return errors.New("Failed to update position")
	}

	if err := tx.Model(&adjacentItem).Update("position", currentPos).Error; err != nil {
		tx.Rollback()
		return errors.New("Failed to update position")
	}

	tx.Commit()

	return nil
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
