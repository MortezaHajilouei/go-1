package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	FirstName    string
	LastName     string
	createdAtUtc time.Time
}
