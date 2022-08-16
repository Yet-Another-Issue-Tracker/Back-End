package project

import (
	"encoding/json"
	"issue-service/app/issue-api/routes/models"
	"issue-service/internal"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getProjectFromRequestBody(r *http.Request) (models.CreateProjectRequest, error) {
	var requestProject models.CreateProjectRequest
	err := json.NewDecoder(r.Body).Decode(&requestProject)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error reading request body")
		return models.CreateProjectRequest{}, &models.ErrorResponse{
			ErrorMessage: "Error reading request body",
			ErrorCode:    400,
		}
	}
	return requestProject, nil
}

func createAddProjectHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		requestBody, err := getProjectFromRequestBody(r)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		validationErr := internal.ValidateRequest(requestBody)
		if validationErr != nil {
			internal.LogAndReturnErrorResponse(validationErr, w)
			return
		}

		requestProject := models.Project{
			Name:   requestBody.Name,
			Client: requestBody.Client,
			Type:   requestBody.Type,
		}

		projectId, err := createProject(
			database,
			requestProject,
		)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		response, err := internal.GetCreateResponseBody(projectId)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}

		w.Write(response)
	}
}

func createGetProjectsHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		projects, err := getProjects(database)
		if err != nil {
			internal.LogAndReturnErrorResponse(err, w)
			return
		}
		response, err := json.Marshal(projects)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling the response")
			internal.LogAndReturnErrorResponse(&models.ErrorResponse{
				ErrorMessage: "Error mashaling the response",
				ErrorCode:    500,
			}, w)
			return
		}
		w.Write(response)
	}
}
