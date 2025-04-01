package controllers

import (
	"go_back/internal/initializers"
	"go_back/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCalendar возвращает список элементов календаря для конкретного пользователя
func GetCalendar(c *fiber.Ctx) error {
	userId := c.Params("userId") // Извлекаем userId из параметров запроса
	var calendarItems []models.CalendarItem
	result := initializers.DB.Where("user_id = ?", userId).Find(&calendarItems)
	if result.Error != nil {
		log.Printf("GetCalendar: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": calendarItems})
}

// NewCalendar создает новый элемент календаря
func NewCalendar(c *fiber.Ctx) error {
	userIdStr := c.Params("userId")      // Извлекаем userId из параметров запроса
	userId, err := uuid.Parse(userIdStr) // Преобразуем строку в uuid.UUID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid userId"})
	}

	var payload models.CreateCalendarItemSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("NewCalendar: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Проверяем, что userId присутствует в payload
	if userId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "userId is required"})
	}

	newCalendarItem := models.CalendarItem{
		ID:        payload.ID,
		UserId:    userId, // Устанавливаем userId из payload
		Date:      payload.Date,
		Title:     payload.Title,
		StartTime: payload.StartTime,
	}

	result := initializers.DB.Create(&newCalendarItem)
	if result.Error != nil {
		log.Printf("NewCalendar: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": newCalendarItem})
}

// UpdateCalendar обновляет элемент календаря
func UpdateCalendar(c *fiber.Ctx) error {
	userIdStr := c.Params("userId")      // Извлекаем userId из параметров запроса
	userId, err := uuid.Parse(userIdStr) // Преобразуем строку в uuid.UUID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid userId"})
	}

	id := c.Params("id")
	var payload models.UpdateCalendarItemSchema

	if err := c.BodyParser(&payload); err != nil {
		log.Printf("UpdateCalendar: BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Проверяем, что userId присутствует в payload
	if userId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "userId is required"})
	}

	var calendarItem models.CalendarItem
	result := initializers.DB.First(&calendarItem, "id = ? AND user_id = ?", id, userId)
	if result.Error != nil {
		log.Printf("UpdateCalendar: CalendarItem not found with ID: %s for userId: %s", id, userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "CalendarItem not found"})
	}

	// Обновляем поля
	if payload.Date != "" {
		calendarItem.Date = payload.Date
	}
	if payload.Title != "" {
		calendarItem.Title = payload.Title
	}
	if payload.StartTime != "" {
		calendarItem.StartTime = payload.StartTime
	}

	initializers.DB.Save(&calendarItem)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": calendarItem})
}

// DeleteCalendar удаляет элемент календаря
func DeleteCalendar(c *fiber.Ctx) error {
	userId := c.Params("userId") // Извлекаем userId из параметров запроса
	id := c.Params("id")

	result := initializers.DB.Delete(&models.CalendarItem{}, "id = ? AND user_id = ?", id, userId)
	if result.RowsAffected == 0 {
		log.Printf("DeleteCalendar: No CalendarItem found with ID: %s for userId: %s", id, userId)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No CalendarItem with that ID exists"})
	} else if result.Error != nil {
		log.Printf("DeleteCalendar: Database error: %v", result.Error)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
