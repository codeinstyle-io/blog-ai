package repository

import (
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/system"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) models.UserRepository {

	if db == nil {
		panic("db is nil")
	}

	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return r.db.Delete(&models.User{}, user).Error
}

func (r *userRepository) FindBySessionToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where(system.CookieName+" = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]*models.User, error) {
	var users []*models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CountByEmail(email string) (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count, err
}

func (r *userRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
