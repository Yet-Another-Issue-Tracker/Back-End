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
	PrepareDatabase(database)

	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return
	}

	testCase.Run("createProject return the new id", func(t *testing.T) {
		database.AutoMigrate(&Project{})

		response, err := createProject(database)

		var foundProject Project

		database.First(&foundProject)

		log.Printf("Project %s", foundProject.Id)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), string(response))
	})
}
