package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/model"
	"github.com/marktrs/simple-todo/repository"

	"github.com/gofiber/fiber/v2"
)

// TaskHandler - handler for task routes
type TaskHandler interface {
	GetAllTasks(c *fiber.Ctx) error
	CreateTask(c *fiber.Ctx) error
	DeleteTask(c *fiber.Ctx) error
	UpdateTask(c *fiber.Ctx) error
}

type taskHandler struct {
	taskRepo  repository.TaskRepository
	validator *validator.Validate
}

func NewTaskHandler(v *validator.Validate, taskRepo repository.TaskRepository) TaskHandler {
	return &taskHandler{
		taskRepo:  taskRepo,
		validator: v,
	}
}

// GetAllTasks query all tasks of a user from the database
func (h *taskHandler) GetAllTasks(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	userId := getUserIDFromToken(token)

	tasks, err := h.taskRepo.GetAllTasks(userId)
	if err != nil {
		return err
	}

	return c.JSON(&model.ListTaskResponse{
		Status: "success",
		Tasks:  tasks,
	})
}

// CreateTask add a new task record to the database
func (h *taskHandler) CreateTask(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)

	var req *model.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.validator.Struct(req); err != nil {
		return err
	}

	var task model.Task
	task.ID = uuid.New().String() // generate a new id for the task
	task.UserId = getUserIDFromToken(token)
	task.Message = req.Message

	if err := h.taskRepo.CreateTask(&task); err != nil {
		return err
	}

	return c.JSON(&model.TaskResponse{
		Status:  "success",
		Message: "task created",
		Task:    &task,
	})
}

// UpdateTask - update a specific task record in the database
func (h *taskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	var req *model.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.validator.Struct(req); err != nil {
		return err
	}

	task, err := h.taskRepo.GetTaskByID(id)
	if err != nil {
		return err
	}

	if task.UserId != getUserIDFromToken(token) {
		return fiber.ErrUnauthorized
	}

	if err = h.validator.Struct(req); err != nil {
		return err
	}

	updates := &model.Task{
		Message:   req.Message,
		Completed: req.Completed,
	}

	if updates.Completed {
		updates.CompletedAt = time.Now()
	}

	task, err = h.taskRepo.UpdateTask(task, updates)
	if err != nil {
		return err
	}

	return c.JSON(&model.TaskResponse{
		Status:  "success",
		Message: "task updated",
		Task:    task,
	})
}

// DeleteTask remove a specific task record from the database
func (h *taskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	var task *model.Task
	var err error

	if task, err = h.taskRepo.GetTaskByID(id); err != nil {
		return err
	}

	if task.UserId != getUserIDFromToken(token) {
		return fiber.ErrUnauthorized
	}

	if err = h.taskRepo.DeleteTask(task); err != nil {
		return err
	}

	return c.JSON(&model.TaskResponse{
		Status:  "success",
		Message: "task deleted",
		Task:    task,
	})
}

func getUserIDFromToken(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string)
}
