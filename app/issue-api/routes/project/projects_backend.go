package project

import (
	"fmt"

	models "issue-service/app/issue-api/routes/models"
	"issue-service/internal"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetProjects(database *gorm.DB) ([]models.Project, error) {
	var projects []models.Project
	result := database.Find(&projects)
	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error retrieving projects")
		return []models.Project{}, &models.ErrorResponse{
			ErrorMessage: "Error retrieving projects",
			ErrorCode:    500,
		}
	}
	return projects, nil
}

func CreateProject(database *gorm.DB, projectName string, projectType string, projectClient string) (uint, error) {
	project := models.Project{Name: projectName, Client: projectClient, Type: projectType}

	result := database.Create(&project)

	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error creating new project")
		if internal.IsDuplicateKeyError(result.Error) {
			return 0, &models.ErrorResponse{
				ErrorMessage: fmt.Sprintf("Project with name \"%s\" already exists", projectName),
				ErrorCode:    409,
			}
		}
		return 0, &models.ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}

	}

	return project.ID, nil
}
