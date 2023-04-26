package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandleHTTPError - custom handle http errors message
func HandleHTTPError(c *fiber.Ctx, err error) error {
	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	body := ErrorResponse{
		Status:  "error",
		Message: fiber.ErrInternalServerError.Message,
	}

	if errors.As(err, &e) {
		body.Message = e.Message
		return c.Status(e.Code).JSON(body)
	}

	// Return status 500 if it's an internal server error
	if err != nil {
		body.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  body.Status,
			"message": body.Message,
		})
	}

	return nil
}
