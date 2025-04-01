package models

import (
	"time"

	"github.com/google/uuid"
)

// CalendarItem представляет модель элемента календаря
type CalendarItem struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserId    uuid.UUID `json:"userId"`
	Date      string    `json:"date"`      // Формат даты, например, "YYYY-MM-DD"
	Title     string    `json:"title"`     // Заголовок записи
	StartTime string    `json:"startTime"` // Время начала
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCalendarItemSchema представляет данные для создания нового элемента календаря
type CreateCalendarItemSchema struct {
	ID        int       `json:"id" binding:"required"`
	UserId    uuid.UUID `json:"userId" binding:"required"`
	Date      string    `json:"date" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	StartTime string    `json:"startTime" binding:"required"`
}

// UpdateCalendarItemSchema представляет данные для обновления элемента календаря
type UpdateCalendarItemSchema struct {
	ID        int       `json:"id" binding:"omitempty"`
	UserId    uuid.UUID `json:"userId" binding:"omitempty"`
	Date      string    `json:"date" binding:"omitempty"`
	Title     string    `json:"title" binding:"omitempty"`
	StartTime string    `json:"startTime" binding:"omitempty"`
}
