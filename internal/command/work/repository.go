package work

import (
	"gorm.io/gorm"
)

type Work struct {
	gorm.Model
	Task     string
	Duration int
}

type Repository struct {
	Database *gorm.DB
}
