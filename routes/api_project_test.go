package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

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

func TestCreateProjectHandler(testCase *testing.T) {
	config, err := GetConfig("../.env")
	if err != nil {
		log.Fatalf("Error reading env configuration: %s", err.Error())
		return
	}
	testRouter := NewRouter(config)
	database, err := ConnectDatabase(config)
	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return
	}

	projectName := "project-name"

	inputProject := Project{
		Name:   projectName,
		Client: "client-name",
		Type:   "project-type",
	}

	testCase.Run("/projects - ok", func(t *testing.T) {
		SetupAndResetDatabase(database)

		expectedResponse := CreateProjectResponse{
			Id: "1",
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		requestBody, err := json.Marshal(inputProject)

		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		bodyReader := bytes.NewReader(requestBody)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodPost, "/v1/projects", bodyReader)
		require.NoError(t, requestError, "Error creating the /projects request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		var foundProject Project
		database.First(&foundProject)

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, string(expectedJsonReponse), string(body), "The response body should be the expected one")
		require.Equal(t, projectName, foundProject.Name)
	})

	testCase.Run("/projects - ko - missing name", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedResponse := ErrorResponse{
			ErrorMessage: "Validation error, field: Project.Name, tag: required",
			ErrorCode:    400,
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		inputProject := Project{
			Client: "client-name",
			Type:   "project-type",
		}

		requestBody, err := json.Marshal(inputProject)

		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		bodyReader := bytes.NewReader(requestBody)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodPost, "/v1/projects", bodyReader)
		require.NoError(t, requestError, "Error creating the /projects request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusBadRequest, statusCode, "The response statusCode should be 400")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, fmt.Sprintf("%s\n", string(expectedJsonReponse)), string(body), "The response body should be the expected one")
	})

	testCase.Run("/projects - ko - missing name and type", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedResponse := ErrorResponse{
			ErrorMessage: "Validation error, field: Project.Name, tag: required\nValidation error, field: Project.Type, tag: required",
			ErrorCode:    400,
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		inputProject := Project{
			Client: "client-name",
		}

		requestBody, err := json.Marshal(inputProject)

		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		bodyReader := bytes.NewReader(requestBody)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodPost, "/v1/projects", bodyReader)
		require.NoError(t, requestError, "Error creating the /projects request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusBadRequest, statusCode, "The response statusCode should be 400")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, fmt.Sprintf("%s\n", string(expectedJsonReponse)), string(body), "The response body should be the expected one")
	})

}
