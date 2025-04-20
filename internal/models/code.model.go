package models

import (
	"time"

	"github.com/google/uuid"
)

// VerificationRequest представляет модель запроса на подтверждение
type VerificationRequest struct {
	ID               uuid.UUID `gorm:"primaryKey" json:"id"`
	Email            string    `gorm:"unique" json:"email"`
	VerificationCode string    `json:"verification_code"`
	CodeExpiry       time.Time `json:"code_expiry"`
	CodeUsed         bool      `json:"code_used"`
}
