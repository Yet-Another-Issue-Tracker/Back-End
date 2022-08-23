package issue

import (
	"fmt"
	"issue-service/app/issue-api/routes/models"
	"issue-service/internal"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createIssue(database *gorm.DB, issue models.Issue) (uint, error) {
	result := database.Create(&issue)

	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error creating new issue")

		if internal.IsForeignKeyError(result.Error) {
			entity := "Sprint"
			value := issue.SprintID
			if strings.Contains(result.Error.Error(), "project") {
				entity = "Project"
				value = issue.ProjectID
			}

			return 0, &models.ErrorResponse{
				ErrorMessage: fmt.Sprintf("%s with id \"%d\" does not exists", entity, value),
				ErrorCode:    404,
			}
		}

		return 0, &models.ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}

	}

	return issue.ID, nil
}
