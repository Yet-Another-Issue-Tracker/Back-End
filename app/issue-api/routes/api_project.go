package routes

import (
	"encoding/json"
	"fmt"
	"issue-service/app/issue-api/routes/makes/models"
	"issue-service/internal"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getProjectFromRequestBody(r *http.Request) (models.Project, error) {
	var requestProject models.Project
	err := json.NewDecoder(r.Body).Decode(&requestProject)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error reading request body")
		return models.Project{}, &models.ErrorResponse{
			ErrorMessage: "Error reading request body",
			ErrorCode:    400,
		}
	}
	return requestProject, nil
}
func getCreateProjectResponseBody(projectId uint) ([]byte, error) {
	response := models.CreateProjectResponse{Id: fmt.Sprint(projectId)}

	responseBody, err := json.Marshal(response)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error marshaling the response")
		return []byte{}, &models.ErrorResponse{
			ErrorMessage: "Error mashaling the response",
			ErrorCode:    500,
		}
	}
	return responseBody, nil
}

func CreateAddProjectHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		requestProject, err := getProjectFromRequestBody(r)
		if err != nil {
			internal.ReturnErrorResponse(err, w)
			return
		}

		validationErr := internal.ValidateRequest(requestProject)
		if validationErr != nil {
			internal.ReturnErrorResponse(validationErr, w)
			return
		}

		projectId, err := internal.CreateProject(
			database,
			requestProject.Name,
			requestProject.Type,
			requestProject.Client,
		)
		if err != nil {
			internal.ReturnErrorResponse(err, w)
			return
		}

		response, err := getCreateProjectResponseBody(projectId)
		if err != nil {
			internal.ReturnErrorResponse(err, w)
			return
		}

		w.Write(response)
	}
}

func CreateGetProjectsHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		projects, err := internal.GetProjects(database)
		if err != nil {
			internal.ReturnErrorResponse(err, w)
			return
		}
		response, err := json.Marshal(projects)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling the response")
			internal.ReturnErrorResponse(&models.ErrorResponse{
				ErrorMessage: "Error mashaling the response",
				ErrorCode:    500,
			}, w)
			return
		}
		w.Write(response)
	}
}
