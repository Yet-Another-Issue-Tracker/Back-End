package routes

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateProject(testCase *testing.T) {

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
	expectedResponse := CreateProjectResponse{
		Id: "1",
	}
	expectedJsonReponse, _ := json.Marshal(expectedResponse)

	testCase.Run("createProject return the new id", func(t *testing.T) {
		SetupAndResetDatabase(database)

		response, err := createProject(database, "", "", "")

		var foundProject Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), string(response))
	})

	testCase.Run("createProject with specific name and type", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedProjectName := "project-name"
		expectedType := "project-type"
		expectedClient := "project-client"

		createProject(database, expectedProjectName, expectedType, expectedClient)

		var foundProject Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, expectedProjectName, foundProject.Name)
		require.Equal(t, expectedType, foundProject.Type)
		require.Equal(t, expectedClient, foundProject.Client)
	})

	testCase.Run("create two projects", func(t *testing.T) {
		SetupAndResetDatabase(database)

		createProject(database, "project-1", "", "")
		createProject(database, "project-2", "", "")

		var foundProjects []Project

		result := database.Find(&foundProjects)
		log.Printf("number of rows %d", result.RowsAffected)
		require.Equal(t, nil, err)
		require.Equal(t, 2, int(result.RowsAffected))
	})

	testCase.Run("createProject returns error if project with same name already exits", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedError := "ERROR: duplicate key value violates unique constraint \"idx_projects_name\" (SQLSTATE 23505)"

		expectedProjectName := "project-name"
		expectedType := "project-type"
		expectedClient := "project-client"

		_, err1 := createProject(database, expectedProjectName, expectedType, expectedClient)

		require.Equal(t, nil, err1)

		_, err2 := createProject(database, expectedProjectName, expectedType, expectedClient)

		require.Equal(t, expectedError, err2.Error())
	})

}
