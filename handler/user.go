package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/model"
	"github.com/marktrs/simple-todo/repository"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler - interface for user handler methods
type UserHandler interface {
	CreateUser(c *fiber.Ctx) error
}

type userHandler struct {
	userRepo  repository.UserRepository
	validator *validator.Validate
}

// NewUserHandler - create a new user handler
func NewUserHandler(v *validator.Validate, userRepo repository.UserRepository) UserHandler {
	return &userHandler{
		userRepo:  userRepo,
		validator: v,
	}
}

// hashPassword - hash a password from a string using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CreateUser - add a new user record to the database
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	var req *model.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := req.Validate(h.validator); err != nil {
		return err
	}

	var user model.User
	user.ID = uuid.New().String()
	user.Username = req.Username

	hash, err := hashPassword(req.Password)
	if err != nil {
		return fiber.ErrBadRequest
	}

	user.Password = hash
	if err = h.userRepo.CreateUser(&user); err != nil {
		return err
	}

	t, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		return errors.Join(err, ErrGenerateToken)
	}

	return c.JSON(model.CreateUserResponse{
		Status:  "success",
		Message: "Created user",
		Token:   t,
		User:    &user,
	})
}
