package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/model"
	"github.com/marktrs/simple-todo/repository"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler - handler for user routes
type UserHandler interface {
	CreateUser(c *fiber.Ctx) error
}

type userHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) UserHandler {
	return &userHandler{
		userRepo: userRepo,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CreateUser - add a new user record to the database
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	var user *model.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	user.ID = uuid.NewString()

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	}

	user.Password = hash
	if err = h.userRepo.CreateUser(user); err != nil {
		message := "Couldn't create user"
		httpCode := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			message = "This username is already exists"
			httpCode = http.StatusConflict
		}
		return c.Status(httpCode).JSON(fiber.Map{"status": "error", "message": message})
	}

	t, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't generate token", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "token": t, "user": user})
}
