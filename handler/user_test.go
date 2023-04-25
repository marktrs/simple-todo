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
	"github.com/marktrs/simple-todo/router"
	repoMock "github.com/marktrs/simple-todo/testutil/mocks/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserHandlerTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	app      *fiber.App
	userRepo *repoMock.MockUserRepository
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.app = fiber.New()
	suite.userRepo = repoMock.NewMockUserRepository(suite.ctrl)
	router.SetupRoutes(suite.app, suite.userRepo, repoMock.NewMockTaskRepository(suite.ctrl))
}

func (suite *UserHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *UserHandlerTestSuite) TestCreateUser() {
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
			description: "create user",
			route:       "/api/users",
			body:        `{"username": "test", "password": "valid_password"}`,
			mock: func() {
				suite.userRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			expectedError: false,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"message":"Created user","status":"success"}`,
		},
		{
			description:   "create user with empty body",
			route:         "/api/users",
			body:          ``,
			mock:          func() {},
			expectedError: false,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"message":"Review your input", "status":"error"}`,
		},
		{
			description:   "create user with a password length = 73 (too long)",
			route:         "/api/users",
			body:          `{"username": "test", "password": "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Praesentium, nihil eveniet?"}`,
			mock:          func() {},
			expectedError: false,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `{"message":"Couldn't hash password", "status":"error"}`,
		},
		{
			description: "create duplicated username",
			route:       "/api/users",
			body:        `{"username": "test", "password": "valid_password"}`,
			mock: func() {
				suite.userRepo.EXPECT().CreateUser(gomock.Any()).Return(gorm.ErrDuplicatedKey)
			},
			expectedError: false,
			expectedCode:  http.StatusConflict,
			expectedBody:  `{"message":"This username is already exists","status":"error"}`,
		},
		{
			description: "create user with db error",
			route:       "/api/users",
			body:        `{"username": "test", "password": "valid_password"}`,
			mock: func() {
				suite.userRepo.EXPECT().CreateUser(gomock.Any()).Return(gorm.ErrInvalidField)
			},
			expectedError: false,
			expectedCode:  http.StatusInternalServerError,
			expectedBody:  `{"message":"Couldn't create user","status":"error"}`,
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
			Message string `json:"message"`
			Status  string `json:"status"`
		}

		var actual, expect ResponseBody

		err = json.Unmarshal(body, &actual)
		suite.NoError(err, test.description)

		err = json.Unmarshal([]byte(test.expectedBody), &expect)
		suite.NoError(err, test.description)

		suite.Assertions.Equal(expect, actual, test.description)
	}

}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
