package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/database"
	"github.com/marktrs/simple-todo/model"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB

	var user model.User
	db.First(&user, id)

	if user.Username == "" {
		return false
	}

	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// CreateUser - add a new user record to the database
func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	user.ID = uuid.NewString()

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})

	}

	user.Password = hash
	if err := db.Create(&user).Error; err != nil {
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
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "token": t, "user": user})
}
