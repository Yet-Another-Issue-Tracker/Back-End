package routes

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateProject(testCase *testing.T) {
	expectedResponse := CreateProjectResponse{
		Id: "1",
	}
	expectedJsonReponse, _ := json.Marshal(expectedResponse)

	config, err := GetConfig("../.env")
	if err != nil {
		log.Fatalf("Error reading env configuration: %s", err.Error())
		return
	}
	database, err := ConnectDatabase(config)

	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return
	}

	testCase.Run("createProject return the new id", func(t *testing.T) {
		SetupAndResetDatabase(database)

		response, err := createProject(database)

		var foundProject Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), string(response))
	})

	testCase.Run("create two projects", func(t *testing.T) {
		SetupAndResetDatabase(database)

		createProject(database)
		createProject(database)

		var foundProjects []Project

		result := database.Find(&foundProjects)
		log.Printf("number of rows %d", result.RowsAffected)
		require.Equal(t, nil, err)
		require.Equal(t, 2, int(result.RowsAffected))
	})
}
