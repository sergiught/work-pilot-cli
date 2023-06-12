package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/sergiught/work-pilot-cli/internal/work"
)

// Connect connects to the sqlite database and returns a DB connection.
func Connect() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open("work-pilot.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = database.AutoMigrate(&work.Task{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}
