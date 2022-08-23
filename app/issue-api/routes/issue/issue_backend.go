package issue

import (
	"issue-service/app/issue-api/routes/models"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createIssue(database *gorm.DB, issue models.Issue) (uint, error) {
	result := database.Create(&issue)

	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error creating new issue")

		return 0, &models.ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}

	}

	return issue.ID, nil
}
