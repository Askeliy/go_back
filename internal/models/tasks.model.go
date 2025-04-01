package models

import (
	"time"

	"github.com/google/uuid"
)

// TaskItem представляет модель задачи
type TaskItem struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserId    uuid.UUID `json:"userId"`
	Title     string    `json:"title"`
	Priority  string    `json:"priority"` // Например, "Неважно"
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTaskItemSchema представляет данные для создания новой задачи
type CreateTaskItemSchema struct {
	ID        int       `json:"id" binding:"required"`
	UserId    uuid.UUID `json:"userId" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Priority  string    `json:"priority" binding:"required"`
	Completed bool      `json:"completed" binding:"required"`
}

// UpdateTaskItemSchema представляет данные для обновления задачи
type UpdateTaskItemSchema struct {
	ID        int       `json:"id" binding:"omitempty"`
	UserId    uuid.UUID `json:"userId" binding:"omitempty"`
	Title     string    `json:"title" binding:"omitempty"`
	Priority  string    `json:"priority" binding:"omitempty"`
	Completed bool      `json:"completed" binding:"omitempty"`
}
