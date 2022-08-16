package project

import (
	"encoding/json"
	"fmt"
	models "issue-service/app/issue-api/routes/models"
	"issue-service/internal"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateProject(testCase *testing.T) {
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

	testCase.Run("createProject return the new id", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		inputProject := models.Project{
			Name:   "",
			Type:   "",
			Client: "",
		}
		response, err := CreateProject(database, inputProject)

		var foundProject models.Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), fmt.Sprint(response))
	})

	testCase.Run("createProject with specific name and type", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		expectedProjectName := internal.GetRandomStringName(10)
		expectedType := "project-type"
		expectedClient := "project-client"
		inputProject := models.Project{
			Name:   expectedProjectName,
			Type:   expectedType,
			Client: expectedClient,
		}
		CreateProject(database, inputProject)

		var foundProject models.Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, expectedProjectName, foundProject.Name)
		require.Equal(t, expectedType, foundProject.Type)
		require.Equal(t, expectedClient, foundProject.Client)
	})

	testCase.Run("create two projects", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)

		inputProject1 := models.Project{
			Name:   internal.GetRandomStringName(10),
			Type:   "",
			Client: "",
		}
		inputProject2 := inputProject1
		inputProject2.Name = internal.GetRandomStringName(10)

		CreateProject(database, inputProject1)
		CreateProject(database, inputProject2)

		var foundProjects []models.Project

		result := database.Find(&foundProjects)
		log.Printf("number of rows %d", result.RowsAffected)
		require.Equal(t, nil, err)
		require.Equal(t, 2, int(result.RowsAffected))
	})

	testCase.Run("createProject returns error if project with same name already exits", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)

		expectedProjectName := internal.GetRandomStringName(10)
		expectedError := fmt.Sprintf("Project with name \"%s\" already exists", expectedProjectName)

		expectedType := "project-type"
		expectedClient := "project-client"
		inputProject := models.Project{
			Name:   expectedProjectName,
			Type:   expectedType,
			Client: expectedClient,
		}

		_, err1 := CreateProject(database, inputProject)

		require.Equal(t, nil, err1)

		_, err2 := CreateProject(database, inputProject)

		require.Equal(t, expectedError, err2.Error())
	})
}

func TestGetProjects(testCase *testing.T) {
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

	testCase.Run("getProjects return a list of projects", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		expectedProjectName := internal.GetRandomStringName(10)
		expectedType := "project-type"
		expectedClient := "project-client"
		inputProject := models.Project{
			Name:   expectedProjectName,
			Type:   expectedType,
			Client: expectedClient,
		}
		CreateProject(database, inputProject)

		expectedResponse := []models.Project{
			{
				Name:   expectedProjectName,
				Type:   expectedType,
				Client: expectedClient,
			},
		}

		foundProjects, err := getProjects(database)

		require.Equal(t, nil, err)
		require.Equal(t, expectedResponse[0].Name, foundProjects[0].Name)
		require.Equal(t, expectedResponse[0].Type, foundProjects[0].Type)
		require.Equal(t, expectedResponse[0].Client, foundProjects[0].Client)
		require.Equal(t, uint(1), foundProjects[0].ID)
	})
}
