package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createProject(database *gorm.DB, projectName string, projectType string, projectClient string) ([]byte, error) {
	project := Project{Name: projectName, Client: projectClient, Type: projectType}

	result := database.Create(&project)

	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error creating new project")
		if IsDuplicateKeyError(result.Error) {
			return nil, &ErrorResponse{
				ErrorMessage: fmt.Sprintf("Project with name \"%s\" already exists", projectName),
				ErrorCode:    409,
			}
		}
		return nil, &ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}

	}

	response := CreateProjectResponse{Id: fmt.Sprint(project.ID)}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return nil, &ErrorResponse{
			ErrorMessage: "Error marshaling the response",
			ErrorCode:    500,
		}
	}

	return responseBody, nil
}

func CreateAddProjectHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var requestProject Project
		err := json.NewDecoder(r.Body).Decode(&requestProject)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error reading request body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		validationError := ValidateRequest(requestProject)

		if validationError != "" {
			response := ErrorResponse{
				ErrorMessage: validationError,
				ErrorCode:    400,
			}
			jsonResponse, _ := json.Marshal(response)

			http.Error(w, string(jsonResponse), http.StatusBadRequest)
			return
		}

		projectId, err := createProject(
			database,
			requestProject.Name,
			requestProject.Type,
			requestProject.Client,
		)

		if err != nil {
			var errorResponse *ErrorResponse
			errors.As(err, &errorResponse)

			jsonResponse, _ := json.Marshal(err)

			http.Error(w, string(jsonResponse), errorResponse.ErrorCode)
			return
		}

		w.Write(projectId)
	}
}
