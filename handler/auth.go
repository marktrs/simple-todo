package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/marktrs/simple-todo/config"
	"github.com/marktrs/simple-todo/repository"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler - handler for auth routes
type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	userRepo repository.UserRepository
}

func NewAuthHandler(userRepo repository.UserRepository) AuthHandler {
	return &authHandler{
		userRepo,
	}
}

// checkPasswordHash - compare password with hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken - generate a jwt token with user id and username as claims
func GenerateToken(id, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

// Login - compare user and password and return an access token
func (h *authHandler) Login(c *fiber.Ctx) error {
	type UserData struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var ud UserData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	username := input.Username
	pass := input.Password

	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User not found"})
		}
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on username", "data": err})
	}

	ud = UserData{
		ID:       user.ID,
		Username: user.Username,
	}

	if !checkPasswordHash(pass, user.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	t, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "token": t, "user": ud})
}
