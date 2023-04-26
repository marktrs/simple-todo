package repository

import (
	"errors"
	"time"

	"github.com/marktrs/simple-todo/database"
	"github.com/marktrs/simple-todo/model"
	"gorm.io/gorm"
)

// TaskRepository - interface for task repository
type TaskRepository interface {
	GetAllTasks(userId string) ([]model.Task, error)
	GetTaskByID(taskId string) (*model.Task, error)
	CreateTask(task *model.Task) error
	DeleteTask(task *model.Task) error
	UpdateTask(task, updates *model.Task) (*model.Task, error)
}

type taskRepository struct{}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

// GetAllTasks query all tasks of a user from the database
func (r *taskRepository) GetAllTasks(userId string) ([]model.Task, error) {
	db := database.DB
	var tasks []model.Task

	if err := db.Where("user_id = ?", userId).Order("created_at desc").Find(&tasks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks, nil
		}
		return nil, err
	}

	return tasks, nil
}

// GetTaskByID query a task by id from the database
func (r *taskRepository) GetTaskByID(taskId string) (*model.Task, error) {
	db := database.DB
	var task *model.Task

	if err := db.Where("id = ?", taskId).Find(&task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// CreateTask add a new task record to the database
func (r *taskRepository) CreateTask(task *model.Task) error {
	db := database.DB
	return db.Create(&task).Error
}

// DeleteTask - delete a specific task record in the database
func (r *taskRepository) DeleteTask(task *model.Task) error {
	db := database.DB
	return db.Delete(&task).Error
}

// UpdateTask - update a specific task record in the database
func (r *taskRepository) UpdateTask(task, updates *model.Task) (*model.Task, error) {
	db := database.DB

	if err := db.Model(&task).Updates(map[string]interface{}{
		"message":    updates.Message,
		"completed":  updates.Completed,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return nil, err
	}

	if err := db.Save(&task).Error; err != nil {
		return nil, err
	}

	return task, nil
}
