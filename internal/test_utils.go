package internal

import (
	"encoding/json"
	models "issue-service/app/issue-api/routes/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func AssertProjectsEquality(t *testing.T, expected []byte, actual []byte) {
	var expectedProjects []models.Project
	json.Unmarshal(expected, &expectedProjects)

	var actualProjects []models.Project
	json.Unmarshal(actual, &actualProjects)

	for index, expectedProject := range expectedProjects {
		require.Equal(t, expectedProject.Name, actualProjects[index].Name, "The response Name should be the expected one")
		require.Equal(t, expectedProject.Type, actualProjects[index].Type, "The response Type should be the expected one")
		require.Equal(t, expectedProject.Client, actualProjects[index].Client, "The response Client should be the expected one")
		require.Equal(t, expectedProject.ID, actualProjects[index].ID, "The response ID should be the expected one")
	}

}

func AssertSprintsEquality(t *testing.T, expected []byte, actual []byte) {
	var expectedSprints []models.GetSprintResponse
	json.Unmarshal(expected, &expectedSprints)

	var actualSprints []models.GetSprintResponse
	json.Unmarshal(actual, &actualSprints)

	for index, expectedSprint := range expectedSprints {
		require.Equal(t, expectedSprint.Number, actualSprints[index].Number, "The response Number should be the expected one")
		require.Equal(t, expectedSprint.Completed, actualSprints[index].Completed, "The response Completed should be the expected one")
		require.Equal(t, expectedSprint.ProjectID, actualSprints[index].ProjectID, "The response ProjectID should be the expected one")
		require.Equal(t, expectedSprint.ID, actualSprints[index].ID, "The response ID should be the expected one")
	}
}

func AssertIssuesEquality(t *testing.T, expected []byte, actual []byte) {
	var expectedIssues []models.GetIssueResponse
	json.Unmarshal(expected, &expectedIssues)

	var actualIssues []models.GetIssueResponse
	json.Unmarshal(actual, &actualIssues)

	for index, expectedIssue := range expectedIssues {
		require.Equal(t, expectedIssue.Title, actualIssues[index].Title, "The response Title should be the expected one")
		require.Equal(t, expectedIssue.Description, actualIssues[index].Description, "The response Description should be the expected one")
		require.Equal(t, expectedIssue.ProjectID, actualIssues[index].ProjectID, "The response ProjectID should be the expected one")
		require.Equal(t, expectedIssue.SprintID, actualIssues[index].SprintID, "The response SprintID should be the expected one")
		require.Equal(t, expectedIssue.ID, actualIssues[index].ID, "The response ID should be the expected one")
	}
}
