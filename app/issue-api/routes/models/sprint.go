package models

import (
	"time"

	"gorm.io/gorm"
)

type Sprint struct {
	gorm.Model
	ID                uint   `gorm:"primaryKey"`
	Number            string `gorm:"uniqueIndex:idx_member"`
	ProjectID         int    `gorm:"uniqueIndex:idx_member"`
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
	Completed         bool      `json:"completed,omitempty"`
	MaxIssuePerSprint int       `json:"maxIssuePerSprint,omitempty"`
}

type CreatePatchRequest struct {
	ID                uint      `json:"id" validate:"required"`
	ProjectId         int       `json:"projectId" validate:"required"`
	Number            string    `json:"number,omitempty"`
	StartDate         time.Time `json:"startDate,omitempty"`
	EndDate           time.Time `json:"endDate,omitempty"`
	Completed         bool      `json:"completed,omitempty"`
	MaxIssuePerSprint int       `json:"maxIssuePerSprint,omitempty"`
}
