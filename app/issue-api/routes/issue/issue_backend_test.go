package issue

import (
	"encoding/json"
	"fmt"
	"issue-service/app/issue-api/routes/models"
	"issue-service/internal"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateIssue(testCase *testing.T) {
	config, err := internal.GetConfig("../../../../.env")
	if err != nil {
		log.Fatalf("Error reading env configuration: %s", err.Error())
		return
	}

	database, err := internal.ConnectDatabase(config)
	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return
	}
	expectedResponse := 1
	expectedJsonReponse, _ := json.Marshal(expectedResponse)
	sprintNumber := "12345"
	expectedTitle := "Task title"
	expectedDescription := "Task description"
	inputIssue := models.Issue{
		Type:        "Task",
		Title:       expectedTitle,
		Description: expectedDescription,
		Status:      "To Do",
		Assignee:    "Assignee",
	}

	testCase.Run("createIssue return the new id", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		projectId, sprintId := internal.CreateProjectAndSprint(database, sprintNumber)
		inputIssue.ProjectID = int(projectId)
		inputIssue.SprintID = int(sprintId)

		response, err := createIssue(database, inputIssue)

		var foundIssue models.Issue

		database.First(&foundIssue)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), fmt.Sprint(response))
		require.Equal(t, int(projectId), foundIssue.ProjectID)
		require.Equal(t, int(sprintId), foundIssue.SprintID)
		require.Equal(t, expectedTitle, foundIssue.Title)
		require.Equal(t, expectedDescription, foundIssue.Description)
	})

	testCase.Run("successfully create two issue on same sprint", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		projectId, sprintId := internal.CreateProjectAndSprint(database, sprintNumber)
		inputIssue.ProjectID = int(projectId)
		inputIssue.SprintID = int(sprintId)

		_, err1 := createIssue(database, inputIssue)
		require.Equal(t, nil, err1)
		_, err2 := createIssue(database, inputIssue)
		require.Equal(t, nil, err2)

		var foundIssue []models.Issue

		database.Find(&foundIssue)

		require.Equal(t, 2, len(foundIssue))
	})

	testCase.Run("createIssue returns error if sprint does not exists", func(t *testing.T) {
		wrongSprintId := 99999
		expectedError := fmt.Sprintf("Sprint with id \"%d\" does not exists", wrongSprintId)

		inputIssue.SprintID = wrongSprintId

		internal.SetupAndResetDatabase(database)
		internal.CreateProjectAndSprint(database, sprintNumber)

		_, err := createIssue(database, inputIssue)

		require.Equal(t, expectedError, err.Error())
	})

	testCase.Run("createIssue returns error if project does not exists", func(t *testing.T) {
		wrongProjectId := 99999
		expectedError := fmt.Sprintf("Project with id \"%d\" does not exists", wrongProjectId)

		inputIssue.ProjectID = wrongProjectId

		internal.SetupAndResetDatabase(database)

		_, err := createIssue(database, inputIssue)

		require.Equal(t, expectedError, err.Error())
	})
}
