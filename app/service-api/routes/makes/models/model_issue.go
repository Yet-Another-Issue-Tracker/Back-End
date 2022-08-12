package models

import (
	"time"

	"gorm.io/gorm"
)

type Issue struct {
	ID          uint   `gorm:"primaryKey"`
	Id          string `json:"id,omitempty"`
	Type_       string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Assignee    string `json:"assignee,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
