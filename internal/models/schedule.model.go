package models

import (
	"time"

	"github.com/google/uuid"
)

// ScheduleItem представляет модель элемента расписания
type ScheduleItem struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserId    uuid.UUID `json:"userId"`
	TimeStart string    `json:"timeStart"`
	TimeEnd   string    `json:"timeEnd"`
	Title     string    `json:"title"`
	Day       int       `json:"day"`
	Type      string    `json:"type"` // Например, "работа"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateScheduleItemSchema представляет данные для создания нового элемента расписания
type CreateScheduleItemSchema struct {
	ID        int       `json:"id" binding:"required"`
	UserId    uuid.UUID `json:"userId" binding:"required"`
	TimeStart string    `json:"timeStart" binding:"required"`
	TimeEnd   string    `json:"timeEnd" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Day       int       `json:"day" binding:"required"`
	Type      string    `json:"type" binding:"required"`
}

// UpdateScheduleItemSchema представляет данные для обновления элемента расписания
type UpdateScheduleItemSchema struct {
	ID        int       `json:"id" binding:"omitempty"`
	UserId    uuid.UUID `json:"userId" binding:"omitempty"`
	TimeStart string    `json:"timeStart" binding:"omitempty"`
	TimeEnd   string    `json:"timeEnd" binding:"omitempty"`
	Title     string    `json:"title" binding:"omitempty"`
	Day       int       `json:"day" binding:"omitempty"`
	Type      string    `json:"type" binding:"omitempty"`
}
