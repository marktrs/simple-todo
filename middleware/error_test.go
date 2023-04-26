package middleware_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/marktrs/simple-todo/middleware"
	"github.com/stretchr/testify/assert"
)

func TestHandleHTTPError(t *testing.T) {
	// Arrange
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "should return 400 bad request",
			err:          fiber.ErrBadRequest,
			expectedCode: 400,
			expectedBody: `{"message":"Bad Request","status":"error"}`,
		},
		{
			name:         "should return 401 unauthorized",
			err:          fiber.ErrUnauthorized,
			expectedCode: 401,
			expectedBody: `{"message":"Unauthorized","status":"error"}`,
		},
		{
			name:         "should return 403 forbidden",
			err:          fiber.ErrForbidden,
			expectedCode: 403,
			expectedBody: `{"message":"Forbidden","status":"error"}`,
		},
		{
			name:         "should return 404 not found",
			err:          fiber.ErrNotFound,
			expectedCode: 404,
			expectedBody: `{"message":"Not Found","status":"error"}`,
		},
		{
			name:         "should return 409 conflict",
			err:          fiber.ErrConflict,
			expectedCode: 409,
			expectedBody: `{"message":"Conflict","status":"error"}`,
		},
		{
			name:         "should return 500 internal server error",
			err:          fiber.ErrInternalServerError,
			expectedCode: 500,
			expectedBody: `{"message":"Internal Server Error","status":"error"}`,
		},
	}

	for _, test := range tests {
		app := fiber.New(
			fiber.Config{
				ErrorHandler: middleware.HandleHTTPError,
			},
		)

		defer func() {
			if err := recover(); err != nil {
				t.Fatalf("Middleware should not panic")
			}
		}()

		app.Get("/", func(c *fiber.Ctx) error {
			return test.err
		})

		req := httptest.NewRequest("GET", "/", nil)

		// Act
		resp, err := app.Test(req, -1)
		assert.NoError(t, err, test.name)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, test.name)

		var actual, expect middleware.ErrorResponse

		err = json.Unmarshal(body, &actual)
		assert.NoError(t, err, test.name)

		err = json.Unmarshal([]byte(test.expectedBody), &expect)
		assert.NoError(t, err, test.name)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, test.expectedCode, resp.StatusCode)
		assert.Equal(t, expect, actual)
	}
}
