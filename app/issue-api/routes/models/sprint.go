package models

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model
	Number            string
	ProjectID         int
	Project           Project
	StartDate         time.Time
	EndDate           time.Time
	Completed         bool
	MaxIssuePerSprint int
}

type CreateSprintRequest struct {
	Number            string    `json:"number,omitempty" validate:"required"`
	StartDate         time.Time `json:"startDate,omitempty"`
	EndDate           time.Time `json:"endDate,omitempty"`
	MaxIssuePerSprint int       `json:"maxIssuePerSprint,omitempty"`
}
