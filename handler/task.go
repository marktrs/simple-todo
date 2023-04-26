package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/model"
	"github.com/marktrs/simple-todo/repository"
	"gorm.io/gorm"

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(fiber.Map{"status": "success", "tasks": tasks})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "couldn't get tasks", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "tasks": tasks})
}

// CreateTask add a new task record to the database
func (h *taskHandler) CreateTask(c *fiber.Ctx) error {
	var task *model.Task

	token := c.Locals("user").(*jwt.Token)

	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "couldn't parse request body to create new task", "data": err})
	}

	// generate a new id for the task
	task.ID = uuid.New().String()
	task.UserId = getUserIDFromToken(token)

	if err := h.taskRepo.CreateTask(task); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "couldn't create a new task", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "created new task", "task": task})
}

// UpdateTask - update a specific task record in the database
func (h *taskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	task, err := h.taskRepo.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "no task found with id"})
		}
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"status": "error", "message": "unable to process operation"})
	}

	if task.UserId != getUserIDFromToken(token) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	var updates *model.Task
	if err = c.BodyParser(&updates); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "bad request"})
	}

	task, err = h.taskRepo.UpdateTask(task, updates)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "couldn't update task", "data": err})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "task successfully updated", "task": task})
}

// DeleteTask remove a specific task record from the database
func (h *taskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	var task *model.Task
	var err error
	if task, err = h.taskRepo.GetTaskByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "no task found with id"})
		}
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"status": "error", "message": "unable to process operation"})
	}

	if task.UserId != getUserIDFromToken(token) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	if err = h.taskRepo.DeleteTask(task); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "couldn't delete task", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "task successfully deleted"})
}

func getUserIDFromToken(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string)
}
