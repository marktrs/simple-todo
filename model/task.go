package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Task struct
type Task struct {
	ID          string         `gorm:"primaryKey,index" json:"id"`
	CreatedAt   time.Time      `gorm:"autoCreateTime,sort:desc" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	Message     string         `gorm:"index,not null" json:"message"`
	Completed   bool           `gorm:"index,not null" json:"completed"`
	CompletedAt time.Time      `gorm:"default:null" json:"completed_at"`
	UserId      string         `gorm:"index,not null" json:"user_id"`
}

type ListTaskResponse struct {
	Status string `json:"status"`
	Tasks  []Task `json:"tasks"`
}

type CreateTaskRequest struct {
	Message string `json:"message" validate:"required,min=1,max=120"`
}

func (r *CreateTaskRequest) Validate() error {
	return validator.New().Struct(r)
}

type UpdateTaskRequest struct {
	Message   string `json:"message,omitempty" validate:"min=1,max=120"`
	Completed bool   `json:"completed,omitempty"`
}

func (r *UpdateTaskRequest) Validate() error {
	return validator.New().Struct(r)
}

type TaskResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Task    *Task  `json:"task"`
}
