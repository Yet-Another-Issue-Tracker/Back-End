package routes

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint   `gorm:"primaryKey"`
	Id        string `json:"id,omitempty"`
	Client    string `json:"client,omitempty"`
	Name      string `gorm:"unique" json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CreateProjectResponse struct {
	Id string `json:"id,omitempty"`
}
