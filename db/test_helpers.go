package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err := ExecuteMigrations(db); err != nil {
		panic(fmt.Sprintf("failed to auto-migrate: %v", err))
	}

	return db
}
