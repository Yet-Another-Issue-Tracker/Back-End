package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint `gorm:"primaryKey"`
	Client    string
	Name      string `gorm:"unique"`
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CreateProjectRequest struct {
	Client string `json:"client,omitempty"`
	Name   string `json:"name,omitempty" validate:"required"`
	Type   string `json:"type,omitempty" validate:"required"`
}
