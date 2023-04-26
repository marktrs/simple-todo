package handler

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/marktrs/simple-todo/config"
	"github.com/marktrs/simple-todo/model"
	"github.com/marktrs/simple-todo/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrGenerateToken = errors.New("error generating token")
)

// AuthHandler - handler for auth routes
type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	userRepo  repository.UserRepository
	validator *validator.Validate
}

func NewAuthHandler(v *validator.Validate, userRepo repository.UserRepository) AuthHandler {
	return &authHandler{
		userRepo:  userRepo,
		validator: v,
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
	var req *model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := req.Validate(h.validator); err != nil {
		return err
	}

	var user *model.User
	var err error

	if user, err = h.userRepo.GetByUsername(req.Username); err != nil {
		return err
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return fiber.ErrUnauthorized
	}

	t, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		return ErrGenerateToken
	}

	return c.JSON(&model.UserResponse{
		Status:  "success",
		Message: "login success",
		Token:   t,
		User:    user,
	})
}
