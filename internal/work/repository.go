package work

import (
	"time"

	"gorm.io/gorm"
)

// Repository used to fetch work
// task data from the database.
type Repository struct {
	Database *gorm.DB
}

// Task is the model stored in the database.
type Task struct {
	gorm.Model
	Name     string
	Duration time.Duration
}

// TableName returns the name of database
// table where work tasks are stored.
func (Task) TableName() string {
	return "work_tasks"
}

// CreateWorkTask will store a work task in the database.
func (r *Repository) CreateWorkTask(task Task) error {
	if task.Name == "" {
		task.Name = "default"
	}

	return r.Database.Create(&task).Error
}

// GetAllWorkTasks will retrieve all work tasks from the database.
func (r *Repository) GetAllWorkTasks() ([]Task, error) {
	var tasks []Task

	if err := r.Database.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
