package database

import (
	"fmt"
	"os"
	"path/filepath"

	gap "github.com/muesli/go-app-paths"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/sergiught/work-pilot-cli/internal/work"
)

// Connect connects to the sqlite database and returns a DB connection.
func Connect() (*gorm.DB, error) {
	scope := gap.NewScope(gap.User, "work-pilot")

	databasePath, err := scope.DataPath("logbook.db")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database scope path: %w", err)
	}

	dir := filepath.Dir(databasePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database dir: %w", err)
		}
	}

	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		file, err := os.Create(databasePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}

		defer func() {
			_ = file.Close()
		}()
	}

	database, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = database.AutoMigrate(&work.Task{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}
