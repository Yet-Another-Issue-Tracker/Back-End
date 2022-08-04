package routes

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint   `gorm:"primaryKey"`
	Id        string `json:"id,omitempty"`
	Client    string `json:"client,omitempty"`
	Name      string `json:"client,omitempty"`
	Type      string `json:"type,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CreateProjectResponse struct {
	Id string `json:"id,omitempty"`
}
