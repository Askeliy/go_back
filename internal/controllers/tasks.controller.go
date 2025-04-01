package controllers

import (
	"go_back/internal/initializers"
	"go_back/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetTasks возвращает список задач для конкретного пользователя
func GetTasks(c *fiber.Ctx) error {
	userId := c.Params("userId") // Извлекаем userId из параметров запроса
	var tasks []models.TaskItem
	result := initializers.DB.Where("user_id = ?", userId).Find(&tasks)
	if result.Error != nil {
		log.Printf("GetTasks: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": tasks})
}

// NewTask создает новую задачу для конкретного пользователя
func NewTask(c *fiber.Ctx) error {
	userIdStr := c.Params("userId")      // Извлекаем userId из параметров запроса
	userId, err := uuid.Parse(userIdStr) // Преобразуем строку в uuid.UUID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid userId"})
	}

	var payload models.CreateTaskItemSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("NewTask: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Проверяем, что userId присутствует в payload
	if userId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "userId is required"})
	}

	newTaskItem := models.TaskItem{
		ID:        payload.ID,
		UserId:    userId, // Устанавливаем userId из payload
		Title:     payload.Title,
		Priority:  payload.Priority,
		Completed: payload.Completed,
	}

	result := initializers.DB.Create(&newTaskItem)
	if result.Error != nil {
		log.Printf("NewTask: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": newTaskItem})
}

// UpdateTask обновляет задачу для конкретного пользователя
func UpdateTask(c *fiber.Ctx) error {
	userIdStr := c.Params("userId")      // Извлекаем userId из параметров запроса
	userId, err := uuid.Parse(userIdStr) // Преобразуем строку в uuid.UUID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid userId"})
	}
	id := c.Params("id")
	var payload models.UpdateTaskItemSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("UpdateTask: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Проверяем, что userId присутствует в payload
	if userId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "userId is required"})
	}

	var taskItem models.TaskItem
	result := initializers.DB.First(&taskItem, "id = ? AND user_id = ?", id, userId)
	if result.Error != nil {
		log.Printf("UpdateTask: TaskItem not found with ID: %s for userId: %s", id, userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "TaskItem not found"})
	}

	// Обновляем поля
	if payload.Title != "" {
		taskItem.Title = payload.Title
	}
	if payload.Priority != "" {
		taskItem.Priority = payload.Priority
	}
	if payload.Completed {
		taskItem.Completed = payload.Completed
	}

	initializers.DB.Save(&taskItem)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": taskItem})
}

// DeleteTask удаляет задачу для конкретного пользователя
func DeleteTask(c *fiber.Ctx) error {
	userId := c.Params("userId") // Извлекаем userId из параметров запроса
	id := c.Params("id")

	result := initializers.DB.Delete(&models.TaskItem{}, "id = ? AND user_id = ?", id, userId)
	if result.RowsAffected == 0 {
		log.Printf("DeleteTask: No TaskItem found with ID: %s for userId: %s", id, userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No TaskItem with that ID exists"})
	} else if result.Error != nil {
		log.Printf("DeleteTask: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
