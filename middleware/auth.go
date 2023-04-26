package middleware

import (
	"github.com/marktrs/simple-todo/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

type AuthErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Protected protect routes
func Protected() fiber.Handler {
	secret := config.Config("SECRET")

	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(secret),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	body := AuthErrorResponse{
		Status:  "error",
		Message: "invalid or expired JWT",
	}

	if err.Error() == "signature is invalid" {
		body.Message = "missing or malformed JWT"
		return c.Status(fiber.StatusBadRequest).JSON(body)
	}

	return c.Status(fiber.StatusUnauthorized).JSON(body)
}
