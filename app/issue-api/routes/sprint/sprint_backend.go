package sprint

import (
	"fmt"
	"issue-service/app/issue-api/routes/models"
	"issue-service/internal"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createSprint(
	database *gorm.DB, sprint models.Sprint) (uint, error) {

	result := database.Create(&sprint)

	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error creating new sprint")
		if internal.IsDuplicateKeyError(result.Error) {
			return 0, &models.ErrorResponse{
				ErrorMessage: fmt.Sprintf("Sprint with number \"%s\" already exists", sprint.Number),
				ErrorCode:    409,
			}
		}

		if internal.IsForeignKeyError(result.Error) {
			return 0, &models.ErrorResponse{
				ErrorMessage: fmt.Sprintf("Project with id \"%d\" does not exists", sprint.ProjectID),
				ErrorCode:    404,
			}
		}
		return 0, &models.ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}

	}

	return sprint.ID, nil
}

func patchSprint(database *gorm.DB, sprint models.Sprint) error {
	result := database.Model(&sprint).Updates(sprint)

	if result.Error != nil {
		if internal.IsForeignKeyError(result.Error) {
			return &models.ErrorResponse{
				ErrorMessage: fmt.Sprintf("Project with id \"%d\" does not exists", sprint.ProjectID),
				ErrorCode:    404,
			}
		}

		return &models.ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}
	}
	if result.RowsAffected == 0 {
		return &models.ErrorResponse{
			ErrorMessage: fmt.Sprintf("Sprint with id \"%d\" does not exists", sprint.ID),
			ErrorCode:    404,
		}
	}
	return nil
}
