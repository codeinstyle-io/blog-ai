package db

import (
	"fmt"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
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
	if err := db.First(&settings).Error; err == gorm.ErrRecordNotFound {
		settings = models.Settings{
			Title:        "Captain",
			Subtitle:     "An AI authored blog engine",
			Timezone:     "UTC",
			ChromaStyle:  "solarized-dark",
			Theme:        "default",
			PostsPerPage: 10,
		}
		if err := db.Create(&settings).Error; err != nil {
			return nil, fmt.Errorf("failed to create default settings: %v", err)
		}
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
