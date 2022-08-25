package models

import (
	"time"

	"gorm.io/gorm"
)

type Issue struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	ProjectID   int  `gorm:"uniqueIndex:idx_member"`
	Project     Project
	SprintID    int `gorm:"uniqueIndex:idx_member"`
	Sprint      Project
	Type        string
	Title       string
	Description string
	Status      string
	Assignee    string
}

type CreateIssueRequest struct {
	Type        string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Assignee    string `json:"assignee,omitempty"`
}

type GetIssueResponse struct {
	ID          uint      `json:"id,omitempty"`
	ProjectID   int       `json:"projectId,omitempty"`
	SprintID    int       `json:"sprintId,omitempty"`
	Type        string    `json:"type,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Assignee    string    `json:"assignee,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
