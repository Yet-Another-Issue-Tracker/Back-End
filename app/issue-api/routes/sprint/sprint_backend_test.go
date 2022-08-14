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
	expectedResponse := 1
	expectedJsonReponse, _ := json.Marshal(expectedResponse)
	inputSprint := models.Sprint{
		Number:    "1",
		ProjectID: 1,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
		Completed: false,
	}

	testCase.Run("createSprint return the new id", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)

		project.CreateProject(database, "project-name", "type", "client")
		response, err := CreateSprint(database, inputSprint)

		var foundSprint models.Sprint

		database.First(&foundSprint)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), fmt.Sprint(response))
	})

	testCase.Run("create two projects", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		project.CreateProject(database, "project-name", "type", "client")

		CreateSprint(database, inputSprint)

		inputSprint2 := inputSprint
		inputSprint2.Number = "2"

		CreateSprint(database, inputSprint2)

		var foundSprints []models.Sprint

		result := database.Find(&foundSprints)
		log.Printf("number of rows %d", result.RowsAffected)
		require.Equal(t, nil, err)
		require.Equal(t, 2, int(result.RowsAffected))
		require.Equal(t, "1", foundSprints[0].Number)
		require.Equal(t, "2", foundSprints[1].Number)
	})
}
