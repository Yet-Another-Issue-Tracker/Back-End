package sprint

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

func getPatchSprintFromRequestBody(r *http.Request) (models.CreatePatchRequest, error) {
	var requestBody models.CreatePatchRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error reading request body")
		return models.CreatePatchRequest{}, &models.ErrorResponse{
			ErrorMessage: "Error reading request body",
			ErrorCode:    400,
		}
	}

	return requestBody, nil
}

func getSprintFromRequestBody(r *http.Request) (models.CreateSprintRequest, error) {
	var requestBody models.CreateSprintRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error reading request body")
		return models.CreateSprintRequest{}, &models.ErrorResponse{
			ErrorMessage: "Error reading request body",
			ErrorCode:    400,
		}
	}

	return requestBody, nil
}

func createAddSprintHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		projectId, ok := vars["projectId"]
		if !ok {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error reading projectId path param from request",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}

		requestBody, err := getSprintFromRequestBody(r)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		validationErr := internal.ValidateRequest(requestBody)
		if validationErr != nil {
			internal.LogAndReturnErrorResponse(validationErr, w)
			return
		}

		requestSprint := models.Sprint{
			Number:            requestBody.Number,
			StartDate:         requestBody.StartDate,
			EndDate:           requestBody.EndDate,
			Completed:         false,
			MaxIssuePerSprint: requestBody.MaxIssuePerSprint,
		}

		requestSprint.ProjectID, err = strconv.Atoi(projectId)
		if err != nil {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error parsing projectId to int",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}

		sprintId, err := createSprint(
			database,
			requestSprint,
		)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		response, err := internal.GetCreateResponseBody(sprintId)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		w.Write(response)
	}
}

func createPatchSprintHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		projectId, ok := vars["projectId"]
		if !ok {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error reading projectId path param from request",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}
		sprintId, ok := vars["sprintId"]
		if !ok {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error reading sprintId path param from request",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}
		requestBody, err := getPatchSprintFromRequestBody(r)
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

		sprintUid, err := strconv.ParseUint(sprintId, 10, 32)
		if err != nil {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error parsing sprintId to uint",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
			return
		}

		requestSprint := models.Sprint{
			ProjectID:         projectIdInt,
			ID:                uint(sprintUid),
			Number:            requestBody.Number,
			StartDate:         requestBody.StartDate,
			EndDate:           requestBody.EndDate,
			Completed:         requestBody.Completed,
			MaxIssuePerSprint: requestBody.MaxIssuePerSprint,
		}

		patchError := patchSprint(database, requestSprint)
		if patchError != nil {
			internal.LogAndReturnErrorResponse(patchError, w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func createGetSprintsHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)

		projectId, ok := vars["projectId"]
		if !ok {
			errorResponse := &models.ErrorResponse{
				ErrorMessage: "Error reading projectId path param from request",
				ErrorCode:    500,
			}
			internal.LogAndReturnErrorResponse(errorResponse, w)
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

		sprints, err := getSprints(database, projectIdInt)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		responseBody, err := json.Marshal(sprints)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		w.Write(responseBody)
	}
}
