package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"issue-service/app/issue-api/routes/models"
	"issue-service/internal"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"
)

func callCreateProjectAPI(createProject models.CreateProjectRequest, testRouter *negroni.Negroni) {
	requestBody, err := json.Marshal(createProject)

	if err != nil {
		log.WithField("error", err.Error()).Error("Error marshaling json")
	}

	bodyReader := bytes.NewReader(requestBody)

	request, _ := http.NewRequest(http.MethodPost, "/v1/projects", bodyReader)
	responseRecorder := httptest.NewRecorder()

	testRouter.ServeHTTP(responseRecorder, request)
}

func callCreateSprintAPI(
	createSprint models.CreateSprintRequest,
	projectId int,
	testRouter *negroni.Negroni,
) {
	requestBody, err := json.Marshal(createSprint)

	if err != nil {
		log.WithField("error", err.Error()).Error("Error marshaling json")
	}

	bodyReader := bytes.NewReader(requestBody)

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/projects/%d/sprints", projectId), bodyReader)
	responseRecorder := httptest.NewRecorder()

	testRouter.ServeHTTP(responseRecorder, request)
}

// TODO remove this and use database.Create()
func callCreateProjectAndSprint(
	testRouter *negroni.Negroni,
) {
	sprintNumber := internal.GetRandomStringName(10)
	inputProject := models.CreateProjectRequest{
		Name:   sprintNumber,
		Type:   "project-type",
		Client: "project-client",
	}
	inputSprint := models.CreateSprintRequest{
		Number:    sprintNumber,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
	}
	callCreateProjectAPI(inputProject, testRouter)
	callCreateSprintAPI(inputSprint, 1, testRouter)
}

// Health routes tests
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

// Projects tests
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

	projectName := internal.GetRandomStringName(10)

	// TODO change this to CreateProjectRequest
	inputProject := models.Project{
		Name:   projectName,
		Client: "client-name",
		Type:   "project-type",
	}

	testCase.Run("/projects - 200 - project created", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)

		expectedResponse := models.CreateResponse{
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
			ErrorMessage: "Validation error, field: CreateProjectRequest.Name, tag: required",
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
			ErrorMessage: "Validation error, field: CreateProjectRequest.Name, tag: required\nValidation error, field: CreateProjectRequest.Type, tag: required",
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
		expectedProjectName := internal.GetRandomStringName(10)
		expectedType := "project-type"
		expectedClient := "project-client"
		inputProject := models.CreateProjectRequest{
			Name:   expectedProjectName,
			Type:   expectedType,
			Client: expectedClient,
		}

		callCreateProjectAPI(inputProject, testRouter)

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

		internal.AssertProjectsEquality(t, expectedJsonReponse, body)
	})
}

// Sprints tests
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

	// TODO change this to CreateSprintRequest
	inputSprint := models.Sprint{
		Number:    sprintNumber,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
	}

	testCase.Run("/sprints - 200 - sprint created", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		inputProject := models.CreateProjectRequest{
			Name:   internal.GetRandomStringName(10),
			Type:   "project-type",
			Client: "project-client",
		}

		callCreateProjectAPI(inputProject, testRouter)

		expectedResponse := models.CreateResponse{
			Id: "1",
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		requestBody, err := json.Marshal(inputSprint)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		bodyReader := bytes.NewReader(requestBody)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodPost, "/v1/projects/1/sprints", bodyReader)
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

func TestPatchSprintHandler(testCase *testing.T) {
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
	projectId := 1
	inputSprint := models.Sprint{
		ProjectID: projectId,
		Completed: false,
		Number:    sprintNumber,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
	}

	testCase.Run("/sprints patch - 204 - sprint patched", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		inputProject := models.Project{
			Name:   internal.GetRandomStringName(10),
			Type:   "project-type",
			Client: "project-client",
		}

		database.Create(&inputProject)
		database.Create(&inputSprint)

		patchSprint := models.PatchSprintRequest{
			ID:        inputSprint.ID,
			Completed: true,
		}
		requestBody, err := json.Marshal(patchSprint)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		responseRecorder := httptest.NewRecorder()
		bodyReader := bytes.NewReader(requestBody)
		request, requestError := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/projects/1/sprints/%d", inputSprint.ID), bodyReader)
		require.NoError(t, requestError, "Error creating the /sprints request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusNoContent, statusCode, "The response statusCode should be 204")

		var foundSprint models.Sprint
		database.First(&foundSprint)
		require.Equal(t, true, foundSprint.Completed)
	})

	testCase.Run("/sprints patch - 404 - project does not exists", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		wrongProjectId := 99999
		inputProject := models.Project{
			Name:   internal.GetRandomStringName(10),
			Type:   "project-type",
			Client: "project-client",
		}

		database.Create(&inputProject)

		database.Create(&inputSprint)
		expectedResponse := models.ErrorResponse{
			ErrorMessage: fmt.Sprintf("Project with id \"%d\" does not exists", wrongProjectId),
			ErrorCode:    404,
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		patchSprint := models.PatchSprintRequest{
			ID:        inputSprint.ID,
			Completed: true,
		}
		requestBody, err := json.Marshal(patchSprint)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		responseRecorder := httptest.NewRecorder()
		bodyReader := bytes.NewReader(requestBody)
		request, requestError := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/projects/%d/sprints/%d", wrongProjectId, inputSprint.ID), bodyReader)
		require.NoError(t, requestError, "Error creating the /sprints request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusNotFound, statusCode, "The response statusCode should be 404")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, fmt.Sprintf("%s\n", string(expectedJsonReponse)), string(body), "The response body should be the expected one")
	})
}

func TestGetSprintsHandler(testCase *testing.T) {
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
	projectId := 1
	inputSprint := models.Sprint{
		ProjectID: projectId,
		Completed: false,
		Number:    sprintNumber,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
	}

	testCase.Run("/sprints get - 200 - sprints returned", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)

		inputProject := models.Project{
			Name:   internal.GetRandomStringName(10),
			Type:   "project-type",
			Client: "project-client",
		}
		newSprintNumber := "98765"
		newSprint := inputSprint
		newSprint.Number = newSprintNumber
		database.Create(&inputProject)
		database.Create(&inputSprint)
		database.Create(&newSprint)

		expectedResponse := []models.GetSprintResponse{
			inputSprint.GetSprintResponseFromSprint(),
			newSprint.GetSprintResponseFromSprint(),
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/projects/%d/sprints", projectId), nil)
		require.NoError(t, requestError, "Error creating the /sprints request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		internal.AssertSprintsEquality(t, expectedJsonReponse, body)
	})
}

// Issue tests
func TestCreateIssueHandler(testCase *testing.T) {
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

	expectedTitle := "Task title"
	expectedDescription := "Task description"
	inputIssue := models.CreateIssueRequest{
		Type:        "Task",
		Title:       expectedTitle,
		Description: expectedDescription,
		Status:      "To Do",
		Assignee:    "Assignee",
	}

	testCase.Run("/issues - 200 - issue created", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		callCreateProjectAndSprint(testRouter)

		expectedResponse := models.CreateResponse{
			Id: "1",
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		requestBody, err := json.Marshal(inputIssue)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		bodyReader := bytes.NewReader(requestBody)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(http.MethodPost, "/v1/projects/1/sprints/1/issues", bodyReader)
		require.NoError(t, requestError, "Error creating the /sprints request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		var foundIssue models.Issue
		database.First(&foundIssue)

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		require.Equal(t, string(expectedJsonReponse), string(body), "The response body should be the expected one")
		require.Equal(t, expectedTitle, foundIssue.Title)
		require.Equal(t, expectedDescription, foundIssue.Description)
	})
}

func TestGetIssuesHandler(testCase *testing.T) {
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

	expectedTitle := "Task title"
	expectedDescription := "Task description"
	projectId := 1
	sprintId := 1
	inputIssue1 := models.Issue{
		Type:        "Task",
		ProjectID:   projectId,
		SprintID:    sprintId,
		Title:       expectedTitle,
		Description: expectedDescription,
		Status:      "To Do",
		Assignee:    "Assignee",
	}
	inputIssue2 := models.Issue{
		Type:        "Task 2",
		ProjectID:   projectId,
		SprintID:    sprintId,
		Title:       expectedTitle,
		Description: expectedDescription,
		Status:      "To Do",
		Assignee:    "Assignee",
	}
	testCase.Run("/issues get - 200 - issues returned", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		callCreateProjectAndSprint(testRouter)
		database.Create(&inputIssue1)
		database.Create(&inputIssue2)

		expectedResponse := []models.GetIssueResponse{
			inputIssue1.GetIssueResponseFromIssue(),
			inputIssue2.GetIssueResponseFromIssue(),
		}

		expectedJsonReponse, _ := json.Marshal(expectedResponse)

		responseRecorder := httptest.NewRecorder()
		request, requestError := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/v1/projects/%d/sprints/%d/issues", projectId, sprintId),
			nil,
		)
		require.NoError(t, requestError, "Error creating the /issues request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusOK, statusCode, "The response statusCode should be 200")

		rawBody := responseRecorder.Result().Body
		body, readBodyError := ioutil.ReadAll(rawBody)
		require.NoError(t, readBodyError)
		internal.AssertIssuesEquality(t, expectedJsonReponse, body)
	})

}

func TestPatchIssueHandler(testCase *testing.T) {
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

	expectedTitle := "Task title"
	expectedDescription := "Task description"
	projectId := 1
	sprintId := 1
	inputIssue := models.Issue{
		Type:        "Task",
		ProjectID:   projectId,
		SprintID:    sprintId,
		Title:       expectedTitle,
		Description: expectedDescription,
		Status:      "To Do",
		Assignee:    "Assignee",
	}

	testCase.Run("/issues patch - 204 - issue is patched", func(t *testing.T) {
		internal.SetupAndResetDatabase(database)
		callCreateProjectAndSprint(testRouter)
		database.Create(&inputIssue)
		expectedStatus := "Completed"
		inputIssue.Status = expectedStatus

		patchIssue := models.PatchIssueRequest{
			ID:     inputIssue.ID,
			Status: expectedStatus,
		}
		requestBody, err := json.Marshal(patchIssue)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling json")
		}

		responseRecorder := httptest.NewRecorder()
		bodyReader := bytes.NewReader(requestBody)

		request, requestError := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/v1/projects/%d/sprints/%d/issues/%d", projectId, sprintId, inputIssue.ID),
			bodyReader,
		)
		require.NoError(t, requestError, "Error creating the /issues request")

		testRouter.ServeHTTP(responseRecorder, request)
		statusCode := responseRecorder.Result().StatusCode
		require.Equal(t, http.StatusNoContent, statusCode, "The response statusCode should be 204")

		var issueSprint models.Issue
		database.First(&issueSprint)
		require.Equal(t, expectedStatus, issueSprint.Status)
		require.Equal(t, expectedTitle, issueSprint.Title)
		require.Equal(t, expectedDescription, issueSprint.Description)
	})
}
