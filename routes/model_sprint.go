package routes

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	ID        uint      `gorm:"primaryKey"`
	Number    string    `json:"number,omitempty"`
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
	Completed bool      `json:"completed,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
