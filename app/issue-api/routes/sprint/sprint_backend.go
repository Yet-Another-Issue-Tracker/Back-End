package sprint

import (
	"issue-service/app/issue-api/routes/models"

	"gorm.io/gorm"
)

// TODO: make this private
func CreateSprint(
	database *gorm.DB, sprint models.Sprint) (uint, error) {

	result := database.Create(&sprint)

	if result.Error != nil {
		return 0, &models.ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}

	}

	return sprint.ID, nil
}
