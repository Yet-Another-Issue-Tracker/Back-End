package sprint

import (
	"encoding/json"
	"fmt"
	"issue-service/app/issue-api/routes/models"
	"issue-service/app/issue-api/routes/project"
	"issue-service/internal"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateSprint(testCase *testing.T) {
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
	expectedSprintNumber := "12345"
	expectedResponse := 1
	expectedJsonReponse, _ := json.Marshal(expectedResponse)
	projectId := 1
	inputSprint := models.Sprint{
		Number:    expectedSprintNumber,
		ProjectID: projectId,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
		Completed: false,
	}

	testCase.Run("createSprint return the new id", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		inputProject := models.Project{
			Name:   internal.GetRandomStringName(10),
			Type:   "project-type",
			Client: "project-client",
		}
		project.CreateProject(database, inputProject)
		response, err := CreateSprint(database, inputSprint)

		var foundSprint models.Sprint

		database.First(&foundSprint)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), fmt.Sprint(response))
	})

	testCase.Run("create two sprint", func(t *testing.T) {
		expectedSprint2Number := "98765"

		internal.SetupAndResetDatabase(database)
		inputProject := models.Project{
			Name:   internal.GetRandomStringName(10),
			Type:   "project-type",
			Client: "project-client",
		}
		project.CreateProject(database, inputProject)

		CreateSprint(database, inputSprint)

		inputSprint2 := inputSprint
		inputSprint2.Number = expectedSprint2Number

		CreateSprint(database, inputSprint2)

		var foundSprints []models.Sprint

		result := database.Find(&foundSprints)
		log.Printf("number of rows %d", result.RowsAffected)
		require.Equal(t, nil, err)
		require.Equal(t, 2, int(result.RowsAffected))
		require.Equal(t, expectedSprintNumber, foundSprints[0].Number)
		require.Equal(t, expectedSprint2Number, foundSprints[1].Number)
	})

	testCase.Run("createSprint returns error if sprint with same number already exits", func(t *testing.T) {
		expectedError := fmt.Sprintf("Sprint with number \"%s\" already exists", expectedSprintNumber)

		internal.SetupAndResetDatabase(database)
		inputProject := models.Project{
			Name:   internal.GetRandomStringName(10),
			Type:   "project-type",
			Client: "project-client",
		}
		project.CreateProject(database, inputProject)

		_, err1 := CreateSprint(database, inputSprint)

		require.Equal(t, nil, err1)

		_, err2 := CreateSprint(database, inputSprint)

		require.Equal(t, expectedError, err2.Error())
	})

	testCase.Run("createSprint returns error if project does not exists", func(t *testing.T) {
		expectedError := fmt.Sprintf("Project with id \"%d\" does not exists", projectId)

		internal.SetupAndResetDatabase(database)

		_, err := CreateSprint(database, inputSprint)

		require.Equal(t, expectedError, err.Error())
	})
}
