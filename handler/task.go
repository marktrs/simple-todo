package handler

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/marktrs/simple-todo/database"
	"github.com/marktrs/simple-todo/model"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

// GetAllTasks query all tasks of a user from the database
func GetAllTasks(c *fiber.Ctx) error {
	db := database.DB
	var tasks []model.Task

	token := c.Locals("user").(*jwt.Token)
	db.Where("user_id = ?", getUserIDFromToken(token)).Order("created_at desc").Find(&tasks)

	return c.JSON(fiber.Map{"status": "success", "tasks": tasks})
}

// CreateTask add a new task record to the database
func CreateTask(c *fiber.Ctx) error {
	db := database.DB
	token := c.Locals("user").(*jwt.Token)
	task := new(model.Task)

	if err := c.BodyParser(task); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "couldn't create a new task", "data": err})
	}

	task.ID = uuid.New().String()
	task.UserId = getUserIDFromToken(token)

	db.Create(&task)
	return c.JSON(fiber.Map{"status": "success", "message": "created new task", "task": task})
}

// UpdateTask - update a specific task record in the database
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	token := c.Locals("user").(*jwt.Token)

	type UpdateTaskInput struct {
		Message   string `json:"message"`
		Completed bool   `json:"completed"`
	}

	var task model.Task
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "no task found with id"})
		}
		return c.Status(402).JSON(fiber.Map{"status": "error", "message": "unable to process operation"})
	}

	if task.UserId != getUserIDFromToken(token) {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	var req UpdateTaskInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "bad request"})
	}

	if err := db.Model(&task).Updates(map[string]interface{}{
		"message":   req.Message,
		"completed": req.Completed,
	}).Error; err != nil {
		return c.Status(402).JSON(fiber.Map{"status": "error", "message": "unable to process operation"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "task successfully updated", "task": task})
}

// DeleteTask remove a specific task record from the database
func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	token := c.Locals("user").(*jwt.Token)

	var task model.Task
	if err := db.Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "no task found with id"})
		}
		return c.Status(402).JSON(fiber.Map{"status": "error", "message": "unable to process operation"})
	}

	if task.UserId != getUserIDFromToken(token) {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}
	db.Delete(&task)
	return c.JSON(fiber.Map{"status": "success", "message": "task successfully deleted"})
}

func getUserIDFromToken(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string)
}
