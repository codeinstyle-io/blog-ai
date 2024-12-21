package db

import (
	"fmt"

	"captain-corp/captain/config"
	"captain-corp/captain/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB.Path), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.GetGormLogLevel()),
	})

	if err != nil {
		return nil, err
	}

	if err := ExecuteMigrations(db); err != nil {
		return nil, err
	}

	// Initialize default settings if they don't exist
	var settings models.Settings
	if err := db.FirstOrCreate(&settings, models.Settings{
		Title:        "Captain",
		Subtitle:     "An AI authored blog engine",
		Timezone:     "UTC",
		ChromaStyle:  "solarized-dark",
		Theme:        "default",
		PostsPerPage: 10,
	}).Error; err != nil {
		return nil, fmt.Errorf("failed to initialize default settings: %w", err)
	}

	return db, nil
}

func ExecuteMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Post{},
		&models.Tag{},
		&models.User{},
		&models.Page{},
		&models.MenuItem{},
		&models.Settings{},
		&models.Media{},
	)
}
