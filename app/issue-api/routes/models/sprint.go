package models

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model
	Number    string `json:"number,omitempty"`
	ProjectID int    `json:"projectId,omitempty"`
	Project   Project
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
	Completed bool      `json:"completed,omitempty"`
}
