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

type PatchSprintRequest struct {
	ID                uint      `json:"id" validate:"required"`
	ProjectID         int       `json:"projectId" validate:"required"`
	Number            string    `json:"number,omitempty"`
	StartDate         time.Time `json:"startDate,omitempty"`
	EndDate           time.Time `json:"endDate,omitempty"`
	Completed         bool      `json:"completed,omitempty"`
	MaxIssuePerSprint int       `json:"maxIssuePerSprint,omitempty"`
}

type GetSprintResponse struct {
	ID                uint      `json:"id"`
	ProjectID         int       `json:"projectId"`
	Number            string    `json:"number,omitempty"`
	StartDate         time.Time `json:"startDate,omitempty"`
	EndDate           time.Time `json:"endDate,omitempty"`
	Completed         bool      `json:"completed"`
	MaxIssuePerSprint int       `json:"maxIssuePerSprint,omitempty"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty"`
}

func (sprint Sprint) GetSprintResponseFromSprint() GetSprintResponse {
	return GetSprintResponse{
		ID:                sprint.ID,
		ProjectID:         sprint.ProjectID,
		Number:            sprint.Number,
		StartDate:         sprint.StartDate,
		EndDate:           sprint.EndDate,
		Completed:         sprint.Completed,
		MaxIssuePerSprint: sprint.MaxIssuePerSprint,
		CreatedAt:         sprint.CreatedAt,
		UpdatedAt:         sprint.UpdatedAt,
	}
}
