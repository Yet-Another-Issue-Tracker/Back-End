package internal

import (
	"encoding/json"
	"fmt"
	models "issue-service/app/service-api/routes/makes/models"
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
	expectedResponse := 1
	expectedJsonReponse, _ := json.Marshal(expectedResponse)

	testCase.Run("createProject return the new id", func(t *testing.T) {
		SetupAndResetDatabase(database)

		response, err := CreateProject(database, "", "", "")

		var foundProject models.Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), fmt.Sprint(response))
	})

	testCase.Run("createProject with specific name and type", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedProjectName := "project-name"
		expectedType := "project-type"
		expectedClient := "project-client"

		CreateProject(database, expectedProjectName, expectedType, expectedClient)

		var foundProject models.Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, expectedProjectName, foundProject.Name)
		require.Equal(t, expectedType, foundProject.Type)
		require.Equal(t, expectedClient, foundProject.Client)
	})

	testCase.Run("create two projects", func(t *testing.T) {
		SetupAndResetDatabase(database)

		CreateProject(database, "project-1", "", "")
		CreateProject(database, "project-2", "", "")

		var foundProjects []models.Project

		result := database.Find(&foundProjects)
		log.Printf("number of rows %d", result.RowsAffected)
		require.Equal(t, nil, err)
		require.Equal(t, 2, int(result.RowsAffected))
	})

	testCase.Run("createProject returns error if project with same name already exits", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedError := "Project with name \"project-name\" already exists"

		expectedProjectName := "project-name"
		expectedType := "project-type"
		expectedClient := "project-client"

		_, err1 := CreateProject(database, expectedProjectName, expectedType, expectedClient)

		require.Equal(t, nil, err1)

		_, err2 := CreateProject(database, expectedProjectName, expectedType, expectedClient)

		require.Equal(t, expectedError, err2.Error())
	})
}

func TestGetProjects(testCase *testing.T) {
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

	testCase.Run("getProjects return a list of projects", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedProjectName := "project-name"
		expectedType := "project-type"
		expectedClient := "project-client"

		CreateProject(database, expectedProjectName, expectedType, expectedClient)

		expectedResponse := []models.Project{
			{
				Name:   expectedProjectName,
				Type:   expectedType,
				Client: expectedClient,
			},
		}

		foundProjects, err := GetProjects(database)

		require.Equal(t, nil, err)
		require.Equal(t, expectedResponse[0].Name, foundProjects[0].Name)
		require.Equal(t, expectedResponse[0].Type, foundProjects[0].Type)
		require.Equal(t, expectedResponse[0].Client, foundProjects[0].Client)
		require.Equal(t, uint(1), foundProjects[0].ID)
	})
}
