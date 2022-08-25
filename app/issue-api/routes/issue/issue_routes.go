package issue

import (
	"encoding/json"
	"issue-service/app/issue-api/routes/models"
	"issue-service/internal"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getProjectIdAndSprintIdFromRequest(request *http.Request) (projectId string, sprintId string, err error) {
	vars := mux.Vars(request)
	projectId, projectOk := vars["projectId"]
	if !projectOk {
		return "", "", &models.ErrorResponse{
			ErrorMessage: "Error reading projectId path param from request",
			ErrorCode:    500,
		}
	}
	sprintId, sprintOk := vars["sprintId"]
	if !sprintOk {
		return "", "", &models.ErrorResponse{
			ErrorMessage: "Error reading sprintId path param from request",
			ErrorCode:    500,
		}
	}
	return projectId, sprintId, nil
}

func getIssueFromRequestBody(request *http.Request) (models.CreateIssueRequest, error) {
	var requestBody models.CreateIssueRequest
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error reading request body")
		return models.CreateIssueRequest{}, &models.ErrorResponse{
			ErrorMessage: "Error reading request body",
			ErrorCode:    400,
		}
	}

	return requestBody, nil
}

func createAddIssueHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		projectId, sprintId, err := getProjectIdAndSprintIdFromRequest(r)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		requestBody, err := getIssueFromRequestBody(r)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		validationErr := internal.ValidateRequest(requestBody)
		if validationErr != nil {
			internal.LogAndReturnErrorResponse(validationErr, w)
			return
		}

		requestIssue := models.Issue{
			Type:        requestBody.Type,
			Title:       requestBody.Title,
			Description: requestBody.Description,
			Status:      requestBody.Status,
			Assignee:    requestBody.Assignee,
		}

		requestIssue.ProjectID, err = strconv.Atoi(projectId)
		if err != nil {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error parsing projectId to int",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}

		requestIssue.SprintID, err = strconv.Atoi(sprintId)
		if err != nil {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error parsing sprintId to int",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}

		issueId, err := createIssue(
			database,
			requestIssue,
		)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		response, err := internal.GetCreateResponseBody(issueId)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		w.Write(response)
	}
}

func createPatchIssueHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusNoContent)
	}
}

func createGetIssuesHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		projectId, sprintId, err := getProjectIdAndSprintIdFromRequest(r)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		projectIdInt, err := strconv.Atoi(projectId)
		if err != nil {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error parsing projectId to int",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}

		sprintIdInt, err := strconv.Atoi(sprintId)
		if err != nil {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error parsing sprintId to int",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}

		issues, err := getIssues(database, projectIdInt, sprintIdInt)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		responseBody, err := json.Marshal(issues)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		w.Write(responseBody)
	}
}
