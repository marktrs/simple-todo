package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// User struct
type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt
	Username  string `gorm:"uniqueIndex;not null" json:"username"`
	Password  string `gorm:"not null" json:"password"`
	Tasks     []Task `gorm:"foreignKey:user_id" json:"tasks"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"min=4,max=72,required"`
	Password string `json:"password" validate:"min=4,max=72,required"`
}

func (r *CreateUserRequest) Validate(validator *validator.Validate) error {
	return validator.Struct(r)
}

type CreateUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
	User    *User  `json:"user"`
}
