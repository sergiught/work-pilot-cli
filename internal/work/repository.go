package work

import (
	"gorm.io/gorm"
)

type Repository struct {
	Database *gorm.DB
}

type Task struct {
	gorm.Model
	Name     string
	Duration int
}

func (Task) TableName() string {
	return "work_tasks"
}

func (r *Repository) CreateWorkTask(task Task) error {
	if task.Name == "" {
		task.Name = "default"
	}

	return r.Database.Create(&task).Error
}

func (r *Repository) GetAllWorkTasks() ([]Task, error) {
	var tasks []Task

	if err := r.Database.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
