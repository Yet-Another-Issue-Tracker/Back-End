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
		projectId, sprintId := internal.CreateProjectAndSprint(database)
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
		projectId, sprintId := internal.CreateProjectAndSprint(database)
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
		internal.CreateProjectAndSprint(database)

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

func TestGetIssues(testCase *testing.T) {
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

	testCase.Run("getIssues return one issue", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		projectId, sprintId := internal.CreateProjectAndSprint(database)
		issueId := internal.CreateTestIssue(database, int(projectId), int(sprintId))

		foundIssues, err := getIssues(database, int(projectId), int(sprintId))

		require.Equal(t, nil, err)
		require.Equal(t, issueId, foundIssues[0].ID)
	})

	testCase.Run("getIssues return a list of issues", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		projectId, sprintId := internal.CreateProjectAndSprint(database)
		issueId1 := internal.CreateTestIssue(database, int(projectId), int(sprintId))
		issueId2 := internal.CreateTestIssue(database, int(projectId), int(sprintId))

		foundIssues, err := getIssues(database, int(projectId), int(sprintId))

		require.Equal(t, nil, err)
		require.Equal(t, 2, len(foundIssues))
		require.Equal(t, issueId1, foundIssues[0].ID)
		require.Equal(t, issueId2, foundIssues[1].ID)
	})

	testCase.Run("getIssues return empty list if project does not exist", func(t *testing.T) {
		nonExistingProjectId := 99999
		internal.SetupAndResetDatabase(database)
		projectId, sprintId := internal.CreateProjectAndSprint(database)
		internal.CreateTestIssue(database, int(projectId), int(sprintId))

		foundIssues, err := getIssues(database, nonExistingProjectId, int(sprintId))

		require.Equal(t, nil, err)
		require.Equal(t, []models.GetIssueResponse{}, foundIssues)
	})

	testCase.Run("getIssues return empty list if sprint does not exist", func(t *testing.T) {
		nonExistingSprinttId := 99999
		internal.SetupAndResetDatabase(database)
		projectId, sprintId := internal.CreateProjectAndSprint(database)
		internal.CreateTestIssue(database, int(projectId), int(sprintId))

		foundIssues, err := getIssues(database, int(projectId), nonExistingSprinttId)

		require.Equal(t, nil, err)
		require.Equal(t, []models.GetIssueResponse{}, foundIssues)
	})
}
