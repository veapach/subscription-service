package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	gorm.Model

	ServiceName string     `json:"service_name" gorm:"not null"`
	Price       uint       `json:"price" gorm:"not null"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	StartDate   time.Time  `json:"start_date" gorm:"not null"`
	EndDate     *time.Time `json:"end_date"`
}
