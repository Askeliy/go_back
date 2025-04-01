package controllers

import (
	"go_back/internal/initializers"
	"go_back/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetSchedule возвращает список элементов расписания для конкретного пользователя
func GetSchedule(c *fiber.Ctx) error {
	userId := c.Params("userId") // Извлекаем userId из параметров запроса
	var scheduleItems []models.ScheduleItem
	result := initializers.DB.Where("user_id = ?", userId).Find(&scheduleItems)
	if result.Error != nil {
		log.Printf("GetSchedule: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": scheduleItems})
}

// NewSchedule создает новый элемент расписания
func NewSchedule(c *fiber.Ctx) error {
	userIdStr := c.Params("userId")      // Извлекаем userId из параметров запроса
	userId, err := uuid.Parse(userIdStr) // Преобразуем строку в uuid.UUID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid userId"})
	}

	var payload models.CreateScheduleItemSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("NewSchedule: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Проверяем, что userId присутствует в payload
	if userId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "userId is required"})
	}

	newScheduleItem := models.ScheduleItem{
		ID:        payload.ID,
		UserId:    userId, // Устанавливаем userId из payload
		TimeStart: payload.TimeStart,
		TimeEnd:   payload.TimeEnd,
		Title:     payload.Title,
		Day:       payload.Day,
		Type:      payload.Type,
	}

	log.Printf("NewSchedule: Inserting item with ID: %d, TimeStart: %s, TimeEnd: %s, Title: %s, Day: %d, Type: %s",
		newScheduleItem.ID, newScheduleItem.TimeStart, newScheduleItem.TimeEnd, newScheduleItem.Title, newScheduleItem.Day, newScheduleItem.Type)

	result := initializers.DB.Create(&newScheduleItem)
	if result.Error != nil {
		log.Printf("NewSchedule: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": newScheduleItem})
}

// UpdateSchedule обновляет элемент расписания
func UpdateSchedule(c *fiber.Ctx) error {
	userIdStr := c.Params("userId")      // Извлекаем userId из параметров запроса
	userId, err := uuid.Parse(userIdStr) // Преобразуем строку в uuid.UUID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid userId"})
	}
	id := c.Params("id")
	var payload models.UpdateScheduleItemSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("UpdateSchedule: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Проверяем, что userId присутствует в payload
	if userId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "userId is required"})
	}

	var scheduleItem models.ScheduleItem
	result := initializers.DB.First(&scheduleItem, "id = ? AND user_id = ?", id, userId)
	if result.Error != nil {
		log.Printf("UpdateSchedule: ScheduleItem not found with ID: %s for userId: %s", id, userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "ScheduleItem not found"})
	}

	// Обновляем поля
	if payload.TimeStart != "" {
		scheduleItem.TimeStart = payload.TimeStart
	}
	if payload.TimeEnd != "" {
		scheduleItem.TimeEnd = payload.TimeEnd
	}
	if payload.Title != "" {
		scheduleItem.Title = payload.Title
	}
	if payload.Day >= 0 && payload.Day <= 6 {
		scheduleItem.Day = payload.Day
	}
	if payload.Type != "" {
		scheduleItem.Type = payload.Type
	}

	initializers.DB.Save(&scheduleItem)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": scheduleItem})
}

// DeleteSchedule удаляет элемент расписания
func DeleteSchedule(c *fiber.Ctx) error {
	userId := c.Params("userId") // Изв лекаем userId из параметров запроса
	id := c.Params("id")

	result := initializers.DB.Delete(&models.ScheduleItem{}, "id = ? AND user_id = ?", id, userId)
	if result.RowsAffected == 0 {
		log.Printf("DeleteSchedule: No ScheduleItem found with ID: %s for userId: %s", id, userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No ScheduleItem with that ID exists"})
	} else if result.Error != nil {
		log.Printf("DeleteSchedule: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
