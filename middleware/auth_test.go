package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/marktrs/simple-todo/handler"
	"github.com/marktrs/simple-todo/middleware"
	"github.com/stretchr/testify/assert"
)

var testJWTSigningKey = "secret"

func TestProtected(t *testing.T) {
	// Set environment secret for JWT signing
	assert.NoError(t, os.Setenv("SECRET", testJWTSigningKey))

	token, err := handler.GenerateToken("id", "username")
	assert.NoError(t, err)

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
	assert.NoError(t, os.Setenv("SECRET", testJWTSigningKey))

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
	assert.NoError(t, os.Setenv("SECRET", testJWTSigningKey))

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
	req.Header.Add("Authorization", "Bearer "+"malformed-token")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
