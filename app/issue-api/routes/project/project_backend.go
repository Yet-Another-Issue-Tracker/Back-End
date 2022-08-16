package project

import (
	"fmt"

	models "issue-service/app/issue-api/routes/models"
	"issue-service/internal"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getProjects(database *gorm.DB) ([]models.Project, error) {
	var projects []models.Project
	result := database.Find(&projects)
	if result.Error != nil {
		return []models.Project{}, &models.ErrorResponse{
			ErrorMessage: "Error retrieving projects",
			ErrorCode:    500,
		}
	}
	return projects, nil
}

// TODO: make this private
func CreateProject(database *gorm.DB, project models.Project) (uint, error) {
	result := database.Create(&project)

	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error creating new project")
		if internal.IsDuplicateKeyError(result.Error) {
			return 0, &models.ErrorResponse{
				ErrorMessage: fmt.Sprintf("Project with name \"%s\" already exists", project.Name),
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
