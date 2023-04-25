package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/marktrs/simple-todo/middleware"
	"github.com/stretchr/testify/assert"
)

var testJWTSigningKey = "secret"
var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"

func TestProtected(t *testing.T) {
	// Set environment secret for JWT signing
	os.Setenv("SECRET", testJWTSigningKey)

	app := fiber.New()

	defer func() {
		// Assert
		if err := recover(); err != nil {
			t.Fatalf("Middleware should not panic")
		}
	}()

	app.Use(middleware.Protected())
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestJWTUnauthorized(t *testing.T) {
	// Set environment secret for JWT signing
	os.Setenv("SECRET", testJWTSigningKey)

	app := fiber.New()

	defer func() {
		// Assert
		if err := recover(); err != nil {
			t.Fatalf("Middleware should not panic")
		}
	}()

	app.Use(middleware.Protected())
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/ok", nil)

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMalformed(t *testing.T) {
	// Set environment secret for JWT signing
	os.Setenv("SECRET", testJWTSigningKey)

	app := fiber.New()

	defer func() {
		// Assert
		if err := recover(); err != nil {
			t.Fatalf("Middleware should not panic")
		}
	}()

	app.Use(middleware.Protected())
	app.Get("/ok", func(c *fiber.Ctx) error {
		return jwt.ErrTokenMalformed
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Add("Authorization", "Bearer "+token+"malformed")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
