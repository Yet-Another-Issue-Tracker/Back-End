package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID     uint `gorm:"primaryKey"`
	Client string
	Name   string `gorm:"unique"`
	Type   string
}

type CreateProjectRequest struct {
	Client string `json:"client,omitempty"`
	Name   string `json:"name,omitempty" validate:"required"`
	Type   string `json:"type,omitempty" validate:"required"`
}
