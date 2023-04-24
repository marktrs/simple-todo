package model

import (
	"time"

	"gorm.io/gorm"
)

// Task struct
type Task struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	Message     string         `gorm:"not null" json:"message"`
	Completed   bool           `gorm:"not null" json:"completed"`
	CompletedAt time.Time      `gorm:"default:null" json:"completed_at"`
	UserId      string         `gorm:"not null" json:"user_id"`
}
