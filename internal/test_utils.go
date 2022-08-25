package internal

import (
	models "issue-service/app/issue-api/routes/models"
	"time"

	"gorm.io/gorm"
)

func CreateTestProject(database *gorm.DB) uint {
	inputProject := models.Project{
		Name:   GetRandomStringName(10),
		Type:   "project-type",
		Client: "project-client",
	}
	database.Create(&inputProject)

	return inputProject.ID
}
func CreateTestSprint(database *gorm.DB, sprintNumber string, projectId int) uint {
	inputSprint := models.Sprint{
		Number:    sprintNumber,
		ProjectID: projectId,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
		Completed: false,
	}
	database.Create(&inputSprint)

	return inputSprint.ID
}
func CreateTestIssue(database *gorm.DB, projectId int, sprintId int) uint {
	inputIssue := models.Issue{
		ProjectID:   projectId,
		SprintID:    sprintId,
		Type:        "Task",
		Title:       "Issue title",
		Description: "Issue description",
		Status:      "To Do",
		Assignee:    "Assignee",
	}
	database.Create(&inputIssue)

	return inputIssue.ID
}
func CreateProjectAndSprint(database *gorm.DB) (projectId uint, sprintId uint) {
	projectId = CreateTestProject(database)
	sprintId = CreateTestSprint(database, "12345", int(projectId))

	return
}
