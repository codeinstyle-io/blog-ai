package repository

import (
	"time"

	"codeinstyle.io/captain/models"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) error {

	return r.db.Create(&session).Error
}

func (r *SessionRepository) Count() int64 {
	var count int64
	r.db.Model(&models.Session{}).Count(&count)
	return count
}

func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.Session{}).Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&models.Session{}).Error
}
