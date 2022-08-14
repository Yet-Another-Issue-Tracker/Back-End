package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"issue-service/app/issue-api/routes/models"
	"issue-service/app/issue-api/routes/project"
	"issue-service/internal"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestStatusRoutes(testCase *testing.T) {
	serviceName := "issue-service"
	serviceVersion := "1.0.0"
	config, err := internal.GetConfig("../../../.env")
	if err != nil {
		log.Fatalf("Error reading env configuration: %s", err.Error())
		return
	}
	testRouter := NewRouter(config)

	testCase.Run("/-/healthz - ok", func(t *testing.T) {
		expectedResponse := fmt.Sprintf("{\"status\":\"OK\",\"name\":\"%s\",\"version\":\"%s\"}", serviceName, serviceVersion)
		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodGet, "/-/healthz", nil)
		require.NoError(t, requestError, "Error creating the /-/healthz request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, expectedResponse, string(body), "The response body should be the expected one")
	})

	testCase.Run("/-/ready - ok", func(t *testing.T) {
		expectedResponse := fmt.Sprintf("{\"status\":\"OK\",\"name\":\"%s\",\"version\":\"%s\"}", serviceName, serviceVersion)
		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodGet, "/-/ready", nil)
		require.NoError(t, requestError, "Error creating the /-/ready request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, expectedResponse, string(body), "The response body should be the expected one")
	})
}

func TestCreateProjectHandler(testCase *testing.T) {
	config, err := internal.GetConfig("../../../.env")
	if err != nil {
		log.Fatalf("Error reading env configuration: %s", err.Error())
		return
	}

	testRouter := NewRouter(config)
	database, err := internal.ConnectDatabase(config)
	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return
	}

	projectName := "project-name"

	inputProject := models.Project{
		Name:   projectName,
		Client: "client-name",
		Type:   "project-type",
	}

	testCase.Run("/projects - 200 - project created", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)

		expectedResponse := models.CreateProjectResponse{
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

		var foundProject models.Project
		database.First(&foundProject)

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, string(expectedJsonReponse), string(body), "The response body should be the expected one")
		require.Equal(t, projectName, foundProject.Name)
	})

	testCase.Run("/projects - 400 - request has wrong types", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		expectedResponse := models.ErrorResponse{
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
		internal.SetupAndResetDatabase(database)
		expectedResponse := models.ErrorResponse{
			ErrorMessage: "Validation error, field: Project.Name, tag: required",
			ErrorCode:    400,
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		inputProject := models.Project{
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
		internal.SetupAndResetDatabase(database)
		expectedResponse := models.ErrorResponse{
			ErrorMessage: "Validation error, field: Project.Name, tag: required\nValidation error, field: Project.Type, tag: required",
			ErrorCode:    400,
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		inputProject := models.Project{
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
		internal.SetupAndResetDatabase(database)
		database.Create(&inputProject)

		expectedResponse := models.ErrorResponse{
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

func TestGetProjectsHandler(testCase *testing.T) {
	config, err := internal.GetConfig("../../../.env")
	if err != nil {
		log.Fatalf("Error reading env configuration: %s", err.Error())
		return
	}
	testRouter := NewRouter(config)
	database, err := internal.ConnectDatabase(config)
	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return
	}

	testCase.Run("/projects - 200 - returned list of projects", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		expectedProjectName := "project-name"
		expectedType := "project-type"
		expectedClient := "project-client"

		project.CreateProject(database, expectedProjectName, expectedType, expectedClient)

		expectedResponse := []models.Project{
			{
				ID:     1,
				Name:   expectedProjectName,
				Type:   expectedType,
				Client: expectedClient,
			},
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodGet, "/v1/projects", nil)
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

func TestCreateSprintHandler(testCase *testing.T) {
	config, err := internal.GetConfig("../../../.env")
	if err != nil {
		log.Fatalf("Error reading env configuration: %s", err.Error())
		return
	}

	testRouter := NewRouter(config)
	database, err := internal.ConnectDatabase(config)
	if err != nil {
		log.Fatalf("Error connecting to database %s", err.Error())
		return
	}

	sprintNumber := "12345"

	inputSprint := models.Sprint{
		Number:    sprintNumber,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
	}

	testCase.Run("/projects - 200 - project created", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)

		expectedResponse := models.CreateProjectResponse{
			Id: "1",
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		requestBody, err := json.Marshal(inputSprint)

		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		bodyReader := bytes.NewReader(requestBody)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodPost, "/v1/sprints", bodyReader)
		require.NoError(t, requestError, "Error creating the /sprints request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		var foundSprint models.Sprint
		database.First(&foundSprint)

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, string(expectedJsonReponse), string(body), "The response body should be the expected one")
		require.Equal(t, sprintNumber, foundSprint.Number)
		require.Equal(t, false, foundSprint.Completed)
	})
}

func assertProjectsEquality(t *testing.T, expected []byte, actual []byte) {
	var expectedProjects []models.Project
	json.Unmarshal(expected, &expectedProjects)

	var actualProjects []models.Project
	json.Unmarshal(actual, &actualProjects)

	for index, expectedProject := range expectedProjects {
		require.Equal(t, expectedProject.Name, actualProjects[index].Name)
		require.Equal(t, expectedProject.Type, actualProjects[index].Type)
		require.Equal(t, expectedProject.Client, actualProjects[index].Client)
		require.Equal(t, expectedProject.ID, actualProjects[index].ID)
	}

}
