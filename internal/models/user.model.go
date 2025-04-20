package models

import (
	"time"

	"github.com/google/uuid"
)

// User представляет модель пользователя
type User struct {
	ID                uuid.UUID `gorm:"primaryKey" json:"id"`
	Name              string    `json:"name"`
	Email             string    `gorm:"unique" json:"email"`
	Password          string    `json:"-"`
	VerificationCode  string    `json:"-"`
	PasswordResetCode string    `json:"-"`
	CodeExpiry        time.Time `json:"-"`
	CodeUsed          bool      `json:"-"`
}

// CreateUser Schema представляет данные для создания нового пользователя
type CreateUserSchema struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateUser Schema представляет данные для обновления пользователя
type UpdateUserSchema struct {
	Name     string `json:"name" binding:"omitempty"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

// LoginUser Schema представляет данные для авторизации пользователя
type LoginUserSchema struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordSchema представляет данные для смены пароля
type ChangePasswordSchema struct {
	Email string `json:"email" binding:"required,email"` // Проверка на пустоту и формат email
}

// ConfirmChangePasswordSchema представляет данные для подтверждения смены пароля
type ConfirmChangePasswordSchema struct {
	Email       string `json:"email" binding:"required,email"`        // Проверка на пустоту и формат email
	Code        string `json:"code" binding:"required"`               // Проверка на пустоту
	NewPassword string `json:"new_password" binding:"required,min=6"` // Проверка на пустоту и минимальную длину
}
