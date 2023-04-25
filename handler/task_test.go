package handler_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/marktrs/simple-todo/model"
	"github.com/marktrs/simple-todo/router"
	repoMock "github.com/marktrs/simple-todo/testutil/mocks/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TaskHandlerTestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	app      *fiber.App
	taskRepo *repoMock.MockTaskRepository
	userRepo *repoMock.MockUserRepository
	user     *model.User
}

var testJWTSigningKey = "secret"
var tokenWithClaims = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI3MDAyMDksInVzZXJfaWQiOiJ0ZXN0X2lkIiwidXNlcm5hbWUiOiJ0ZXN0In0.XiOkXAVzZ5hW5DbviA2juzPRVuaBEuq19d1dMlVC2uU"

func (suite *TaskHandlerTestSuite) SetupTest() {
	os.Setenv("SECRET", testJWTSigningKey)

	suite.user = &model.User{
		ID:       "test_id",
		Username: "test",
		Password: "$2a$14$j.2bobj6FGyKREHuYXIqoeSN5TA/Vq1C6dkfzg2zuf3GsGeKFla9u", // pass = valid_password
	}

	suite.ctrl = gomock.NewController(suite.T())
	suite.app = fiber.New()
	suite.taskRepo = repoMock.NewMockTaskRepository(suite.ctrl)
	suite.userRepo = repoMock.NewMockUserRepository(suite.ctrl)
	router.SetupRoutes(suite.app, suite.userRepo, suite.taskRepo)
}

func (suite *TaskHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *TaskHandlerTestSuite) TestGetAllTasks() {
	tasks := []model.Task{
		{
			ID:      "1",
			Message: "test",
			UserId:  "test_id",
		},
	}

	tests := []struct {
		description  string
		requireAuth  bool
		expectedBody string
		mockFunc     func()
		expectedCode int
	}{
		{
			description: "get all tasks with authorized user",
			requireAuth: true,
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetAllTasks(gomock.Any()).Return(tasks, nil)
			},
			expectedBody: `{"status":"success","tasks":[{"id":"1","message":"test","user_id":"test_id"}]}`,
			expectedCode: http.StatusOK,
		},
		// {
		// 	description: "get all tasks with authorized user with db error not found",
		// 	requireAuth: true,
		// 	mockFunc: func() {
		// 		suite.taskRepo.EXPECT().GetAllTasks(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
		// 	},
		// 	expectedBody: `{"status":"success","tasks": nil}`,
		// 	expectedCode: http.StatusOK,
		// },
		{
			description: "get all tasks with authorized user and db error",
			requireAuth: true,
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetAllTasks(gomock.Any()).Return(nil, gorm.ErrInvalidDB)
			},
			expectedBody: `{"status":"error", "message": "couldn't get tasks"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)

		if test.requireAuth {
			req.Header.Add("Authorization", "Bearer "+tokenWithClaims)
		}

		test.mockFunc()

		// Act
		res, err := suite.app.Test(req, -1)

		// Assert
		suite.NoError(err)
		suite.Equal(test.expectedCode, res.StatusCode)

		// Read the response body
		body, err := io.ReadAll(res.Body)
		suite.NoError(err, test.description)

		// Assert the response body
		type ResponseBody struct {
			Status string       `json:"status"`
			Tasks  []model.Task `json:"tasks"`
		}

		var actual, expect ResponseBody

		err = json.Unmarshal(body, &actual)
		suite.NoError(err, test.description)

		err = json.Unmarshal([]byte(test.expectedBody), &expect)
		suite.NoError(err, test.description)

		log.Println(test.expectedBody)

		suite.Assertions.Equal(expect, actual, test.description)
	}
}

func (suite *TaskHandlerTestSuite) TestCreateTask() {
	tests := []struct {
		description string

		mockFunc    func()
		requireAuth bool

		body         string
		expectedBody string
		expectedCode int
	}{
		{
			description: "create a task with authorized user",
			mockFunc: func() {
				suite.taskRepo.EXPECT().CreateTask(gomock.Any()).Return(nil)
			},
			requireAuth:  true,
			body:         `{"message":"test"}`,
			expectedBody: `{"status":"success", "message": "created new task"}`,
			expectedCode: http.StatusOK,
		},
		{
			description:  "create a task with invalid body",
			mockFunc:     func() {},
			requireAuth:  true,
			body:         `{"message": nil}`,
			expectedBody: `{"status":"error", "message": "couldn't parse request body to create new task"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			description: "create a task with db error",
			mockFunc: func() {
				suite.taskRepo.EXPECT().CreateTask(gomock.Any()).Return(gorm.ErrInvalidDB)
			},
			requireAuth:  true,
			body:         `{"message":"test"}`,
			expectedBody: `{"status":"error", "message": "couldn't create a new task"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")
		if test.requireAuth {
			req.Header.Add("Authorization", "Bearer "+tokenWithClaims)
		}

		test.mockFunc()

		// Act
		res, err := suite.app.Test(req, -1)

		// Assert
		suite.NoError(err)
		suite.Equal(test.expectedCode, res.StatusCode)

		// Read the response body
		body, err := io.ReadAll(res.Body)
		suite.NoError(err, test.description)

		// Assert the response body
		type ResponseBody struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		var actual, expect ResponseBody

		err = json.Unmarshal(body, &actual)
		suite.NoError(err, test.description)

		err = json.Unmarshal([]byte(test.expectedBody), &expect)
		suite.NoError(err, test.description)

		suite.Assertions.Equal(expect, actual, test.description)
	}
}

func (suite *TaskHandlerTestSuite) TestUpdateTask() {
	task := model.Task{
		ID:      "1",
		Message: "test",
		UserId:  "test_id",
	}

	tests := []struct {
		description string

		mockFunc    func()
		requireAuth bool

		body         string
		expectedBody string
		expectedCode int
	}{
		{
			description: "update a task with id with authorized user",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(&task, nil)
				suite.taskRepo.EXPECT().UpdateTask(&task, gomock.Any()).Return(&model.Task{
					ID:        task.ID,
					Message:   "new message",
					Completed: true,
				}, nil)
			},
			requireAuth:  true,
			body:         `{"message":"new message", "completed": true}`,
			expectedBody: `{"status":"success", "message": "task successfully updated"}`,
			expectedCode: http.StatusOK,
		},
		{
			description: "update a task with an invalid id",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(nil, gorm.ErrRecordNotFound)
			},
			requireAuth:  true,
			body:         `{"message":"new message", "completed": true}`,
			expectedBody: `{"status":"error", "message": "no task found with id"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			description: "update a task with db error",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(nil, gorm.ErrInvalidDB)
			},
			requireAuth:  true,
			body:         `{"message":"new message", "completed": true}`,
			expectedBody: `{"status":"error", "message": "unable to process operation"}`,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			description: "update a task with an invalid user",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(&model.Task{
					ID:        task.ID,
					Message:   "new message",
					Completed: true,
					UserId:    "another_user",
				}, nil)
			},
			requireAuth:  true,
			body:         `{"message":"new message", "completed": true}`,
			expectedBody: `{"status":"error", "message": "unauthorized"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			description: "update a task with id with authorized user",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(&task, nil)
				suite.taskRepo.EXPECT().UpdateTask(&task, gomock.Any()).Return(nil, gorm.ErrInvalidDB)
			},
			requireAuth:  true,
			body:         `{"message":"new message", "completed": true}`,
			expectedBody: `{"status":"error", "message": "couldn't update task"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodPut, "/api/tasks/"+task.ID, strings.NewReader(test.body))
		req.Header.Set("Content-Type", "application/json")

		if test.requireAuth {
			req.Header.Add("Authorization", "Bearer "+tokenWithClaims)
		}

		test.mockFunc()

		// Act
		res, err := suite.app.Test(req, -1)

		// Assert
		suite.NoError(err)
		suite.Equal(test.expectedCode, res.StatusCode)

		// Read the response body
		body, err := io.ReadAll(res.Body)
		suite.NoError(err, test.description)

		// Assert the response body
		type ResponseBody struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		var actual, expect ResponseBody

		err = json.Unmarshal(body, &actual)
		suite.NoError(err, test.description)

		err = json.Unmarshal([]byte(test.expectedBody), &expect)
		suite.NoError(err, test.description)

		suite.Assertions.Equal(expect, actual, test.description)
	}
}

func (suite *TaskHandlerTestSuite) TestDeleteTask() {
	task := model.Task{
		ID:      "1",
		Message: "test",
		UserId:  "test_id",
	}

	tests := []struct {
		description string

		mockFunc    func()
		requireAuth bool

		expectedBody string
		expectedCode int
	}{
		{
			description: "delete a task with id with authorized user",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(&task, nil)
				suite.taskRepo.EXPECT().DeleteTask(&task).Return(nil)
			},
			requireAuth:  true,
			expectedBody: `{"status":"success", "message": "task successfully deleted"}`,
			expectedCode: http.StatusOK,
		},
		{
			description: "delete a task with invalid id",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(nil, gorm.ErrRecordNotFound)
			},
			requireAuth:  true,
			expectedBody: `{"status":"error", "message": "no task found with id"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			description: "delete a task with db error",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(nil, gorm.ErrInvalidDB)
			},
			requireAuth:  true,
			expectedBody: `{"status":"error", "message": "unable to process operation"}`,
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			description: "delete a task with an invalid user",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(&model.Task{
					ID:        task.ID,
					Message:   "new message",
					Completed: true,
					UserId:    "another_user",
				}, nil)
			},
			requireAuth:  true,
			expectedBody: `{"status":"error", "message": "unauthorized"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			description: "delete a task with id with db error",
			mockFunc: func() {
				suite.taskRepo.EXPECT().GetTaskByID(task.ID).Return(&task, nil)
				suite.taskRepo.EXPECT().DeleteTask(&task).Return(gorm.ErrInvalidDB)
			},
			requireAuth:  true,
			expectedBody: `{"status":"error", "message": "couldn't delete task"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodDelete, "/api/tasks/"+task.ID, nil)

		if test.requireAuth {
			req.Header.Add("Authorization", "Bearer "+tokenWithClaims)
		}

		test.mockFunc()

		// Act
		res, err := suite.app.Test(req, -1)

		// Assert
		suite.NoError(err)
		suite.Equal(test.expectedCode, res.StatusCode)

		// Read the response body
		body, err := io.ReadAll(res.Body)
		suite.NoError(err, test.description)

		// Assert the response body
		type ResponseBody struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		var actual, expect ResponseBody

		err = json.Unmarshal(body, &actual)
		suite.NoError(err, test.description)

		err = json.Unmarshal([]byte(test.expectedBody), &expect)
		suite.NoError(err, test.description)

		suite.Assertions.Equal(expect, actual, test.description)
	}
}

func TestTaskHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskHandlerTestSuite))
}
