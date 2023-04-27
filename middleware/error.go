package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/marktrs/simple-todo/logger"
	"gorm.io/gorm"
)

var (
	validationErrMsg = "Failed input validation"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ErrorResponse struct {
	Status           string            `json:"status"`
	Message          string            `json:"message"`
	ValidationErrors []ValidationError `json:"validation_error,omitempty"`
}

// HandleHTTPError - custom handle http errors message
func HandleHTTPError(c *fiber.Ctx, err error) error {
	header := c.GetReqHeaders()

	rq := c.Request()
	rs := c.Response()

	logger.Log.
		Error().Err(err).
		Str("request_id", c.GetRespHeader("Request-ID")).
		Str("protocol", "http").
		Str("method", string(rq.Header.Method())).
		Str("path", rq.URI().String()).
		Int("status_code", rs.StatusCode()).
		Str("status_text", http.StatusText(rs.StatusCode())).
		Msg("error while handling request")

	startTime, parseErr := time.Parse(time.Layout, header["Start"])

	if parseErr != nil {
		logger.Log.Error().Err(err).
			Str("request_id", c.GetRespHeader("Request-ID")).
			Msg("error while parsing request start time")
	} else {
		duration := time.Since(startTime)
		logger.Log.Error().
			Str("duration", duration.String()).
			Str("request_id", c.GetRespHeader("Request-ID")).
			Msg("request duration")
	}

	body := ErrorResponse{
		Status:  "error",
		Message: fiber.ErrInternalServerError.Message,
	}

	// Evaluate error message format if it's a validation error
	if v, ok := err.(validator.ValidationErrors); ok {
		body.Message = validationErrMsg
		body.ValidationErrors = describeValidationError(v)
		return c.Status(http.StatusBadRequest).JSON(body)
	}

	// If it's a gorm.Error we can handle some the status code and message
	if ok, code := isHandledDBError(&body, err); ok {
		return c.Status(code).JSON(body)
	}

	// If it's a fiber.Error we can retrieve the status code and message
	var e *fiber.Error
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

func describeValidationError(verr validator.ValidationErrors) []ValidationError {
	errs := []ValidationError{}
	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, ValidationError{Field: f.Field(), Reason: err})
	}

	return errs
}

func isHandledDBError(body *ErrorResponse, err error) (bool, int) {
	switch err {
	case gorm.ErrDuplicatedKey:
		body.Message = gorm.ErrDuplicatedKey.Error()
		return true, http.StatusConflict
	case gorm.ErrInvalidField:
		body.Message = gorm.ErrInvalidField.Error()
		return true, http.StatusBadRequest
	case gorm.ErrRecordNotFound:
		body.Message = gorm.ErrRecordNotFound.Error()
		return true, http.StatusNotFound
	default:
		return false, 0
	}
}
