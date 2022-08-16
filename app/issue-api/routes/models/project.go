package models

import (
	"time"

	"gorm.io/gorm"
)

//TODO adapt this to sprint structs
type Project struct {
	ID        uint   `gorm:"primaryKey"`
	Client    string `json:"client,omitempty"`
	Name      string `gorm:"unique" json:"name,omitempty" validate:"required"`
	Type      string `json:"type,omitempty" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
