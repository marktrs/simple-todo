package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/marktrs/simple-todo/middleware"
	"github.com/marktrs/simple-todo/model"
	"github.com/marktrs/simple-todo/router"
	"github.com/marktrs/simple-todo/server"
	repoMock "github.com/marktrs/simple-todo/testutil/mocks/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	app      *fiber.App
	userRepo *repoMock.MockUserRepository
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.app = server.New().App()
	suite.userRepo = repoMock.NewMockUserRepository(suite.ctrl)
	router.SetupRoutes(suite.app, suite.userRepo, repoMock.NewMockTaskRepository(suite.ctrl))
}

func (suite *AuthHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *AuthHandlerTestSuite) TestLogin() {
	actualUser := &model.User{
		ID:       "test_id",
		Username: "test",
		Password: "$2a$14$j.2bobj6FGyKREHuYXIqoeSN5TA/Vq1C6dkfzg2zuf3GsGeKFla9u", // pass = valid_password
	}

	tests := []struct {
		description string

		// Test input
		route  string
		method string
		body   string
		mock   func()

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description: "authenticate user without registration",
			route:       "/api/auth/login",
			body:        `{"username": "test", "password": "test"}`,
			mock: func() {
				suite.userRepo.EXPECT().GetByUsername(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: false,
			expectedCode:  http.StatusNotFound,
			expectedBody:  `{"message":"record not found","status":"error"}`,
		},
		{
			description:   "authenticate user with validation error",
			route:         "/api/auth/login",
			body:          `{}`,
			mock:          func() {},
			expectedError: false,
			expectedCode:  http.StatusBadRequest,
			expectedBody: `{
				"message":"Failed input validation",
				"status":"error",
				"validation_error":[
					{"field":"Username","reason":"min=4"},
					{"field":"Password","reason":"min=4"}
				]}`,
		},
		{
			description: "authenticate user with invalid credentials",
			route:       "/api/auth/login",
			body:        `{"username": "test", "password": "invalid_password"}`,
			mock: func() {
				suite.userRepo.EXPECT().GetByUsername(gomock.Any()).Return(actualUser, nil)
			},
			expectedError: false,
			expectedCode:  http.StatusUnauthorized,
			expectedBody:  `{"message":"Unauthorized", "status":"error"}`,
		},
		{
			description: "authenticate user with valid credentials",
			route:       "/api/auth/login",
			body:        `{"username": "test", "password": "valid_password"}`,
			mock: func() {
				suite.userRepo.EXPECT().GetByUsername(gomock.Any()).Return(actualUser, nil)
			},
			expectedError: false,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"message":"login success", "status":"success"}`,
		},
	}

	for _, test := range tests {
		// Read the response body
		req := httptest.NewRequest(
			fiber.MethodPost,
			test.route,
			strings.NewReader(test.body),
		)

		// Set the content type to JSON
		req.Header.Set("Content-Type", "application/json")

		// Call mock setup if any
		test.mock()

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := suite.app.Test(req, -1)
		suite.Equal(test.expectedCode, res.StatusCode, test.description)

		if test.expectedError {
			suite.Error(err, test.description)
		} else {
			suite.NoError(err, test.description)
		}

		// Read the response body
		body, err := io.ReadAll(res.Body)
		suite.NoError(err, test.description)

		// Assert the response body
		type ResponseBody struct {
			Message          string                       `json:"message"`
			Status           string                       `json:"status"`
			ValidationErrors []middleware.ValidationError `json:"validation_error,omitempty"`
		}

		var actual, expect ResponseBody

		err = json.Unmarshal(body, &actual)
		suite.NoError(err, test.description)

		err = json.Unmarshal([]byte(test.expectedBody), &expect)
		suite.NoError(err, test.description)

		suite.Assertions.Equal(expect, actual, test.description)
	}
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
