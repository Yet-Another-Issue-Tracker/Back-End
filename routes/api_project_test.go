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
	expectedResponse := 1
	expectedJsonReponse, _ := json.Marshal(expectedResponse)

	testCase.Run("createProject return the new id", func(t *testing.T) {
		SetupAndResetDatabase(database)

		response, err := createProject(database, "", "", "")

		var foundProject Project

		database.First(&foundProject)

		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), fmt.Sprint(response))
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
		expectedError := "Project with name \"project-name\" already exists"

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

	testCase.Run("/projects - 200 - project created", func(t *testing.T) {
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

	testCase.Run("/projects - 400 - request has wrong types", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedResponse := ErrorResponse{
			ErrorMessage: "Error reading request body",
			ErrorCode:    400,
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		type WrongProject struct {
			Name bool
		}

		inputProject := WrongProject{
			Name: true,
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
		require.Equal(t, expectedResponse.ErrorCode, statusCode, "The response statusCode should be 500")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, fmt.Sprintf("%s\n", string(expectedJsonReponse)), string(body), "The response body should be the expected one")
	})

	testCase.Run("/projects - 400 - missing name", func(t *testing.T) {
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
		require.Equal(t, expectedResponse.ErrorCode, statusCode, "The response statusCode should be 400")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, fmt.Sprintf("%s\n", string(expectedJsonReponse)), string(body), "The response body should be the expected one")
	})

	testCase.Run("/projects - 400 - missing name and type", func(t *testing.T) {
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
		require.Equal(t, expectedResponse.ErrorCode, statusCode, "The response statusCode should be 400")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, fmt.Sprintf("%s\n", string(expectedJsonReponse)), string(body), "The response body should be the expected one")
	})

	testCase.Run("/projects - 409 - project already exists", func(t *testing.T) {
		SetupAndResetDatabase(database)
		database.Create(&inputProject)

		expectedResponse := ErrorResponse{
			ErrorMessage: fmt.Sprintf("Project with name \"%s\" already exists", projectName),
			ErrorCode:    409,
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
		require.Equal(t, expectedResponse.ErrorCode, statusCode, "The response statusCode should be 409")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, fmt.Sprintf("%s\n", string(expectedJsonReponse)), string(body), "The response body should be the expected one")
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

		createProject(database, expectedProjectName, expectedType, expectedClient)

		expectedResponse := []Project{
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
func assertProjectsEquality(t *testing.T, expected []byte, actual []byte) {
	var expectedProjects []Project
	json.Unmarshal(expected, &expectedProjects)

	var actualProjects []Project
	json.Unmarshal(actual, &actualProjects)

	for index, expectedProject := range expectedProjects {
		require.Equal(t, expectedProject.Name, actualProjects[index].Name)
		require.Equal(t, expectedProject.Type, actualProjects[index].Type)
		require.Equal(t, expectedProject.Client, actualProjects[index].Client)
		require.Equal(t, expectedProject.ID, actualProjects[index].ID)
	}

}
func TestGetProjectsHandler(testCase *testing.T) {
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

	testCase.Run("/projects - 200 - project created", func(t *testing.T) {
		SetupAndResetDatabase(database)
		expectedProjectName := "project-name"
		expectedType := "project-type"
		expectedClient := "project-client"

		createProject(database, expectedProjectName, expectedType, expectedClient)

		expectedResponse := []Project{
			{
				ID:     1,
				Name:   expectedProjectName,
				Type:   expectedType,
				Client: expectedClient,
			},
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		requestBody, err := json.Marshal(inputProject)

		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		bodyReader := bytes.NewReader(requestBody)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodGet, "/v1/projects", bodyReader)
		require.NoError(t, requestError, "Error creating the /projects request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)

		assertProjectsEquality(t, expectedJsonReponse, body)
	})
}
