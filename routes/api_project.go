package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func getProjects(database *gorm.DB) ([]Project, error) {
	var projects []Project
	result := database.Find(&projects)
	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error retrieving projects")
		return []Project{}, &ErrorResponse{
			ErrorMessage: "Error retrieving projects",
			ErrorCode:    500,
		}
	}
	return projects, nil
}
func createProject(database *gorm.DB, projectName string, projectType string, projectClient string) (uint, error) {
	project := Project{Name: projectName, Client: projectClient, Type: projectType}

	result := database.Create(&project)

	if result.Error != nil {
		log.WithField("error", result.Error.Error()).Error("Error creating new project")
		if IsDuplicateKeyError(result.Error) {
			return 0, &ErrorResponse{
				ErrorMessage: fmt.Sprintf("Project with name \"%s\" already exists", projectName),
				ErrorCode:    409,
			}
		}
		return 0, &ErrorResponse{
			ErrorMessage: result.Error.Error(),
			ErrorCode:    500,
		}

	}

	return project.ID, nil
}

func getProjectFromRequestBody(r *http.Request) (Project, error) {
	var requestProject Project
	err := json.NewDecoder(r.Body).Decode(&requestProject)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error reading request body")
		return Project{}, &ErrorResponse{
			ErrorMessage: "Error reading request body",
			ErrorCode:    400,
		}
	}
	return requestProject, nil
}
func getCreateProjectResponseBody(projectId uint) ([]byte, error) {
	response := CreateProjectResponse{Id: fmt.Sprint(projectId)}

	responseBody, err := json.Marshal(response)
	if err != nil {
		log.WithField("error", err.Error()).Error("Error marshaling the response")
		return []byte{}, &ErrorResponse{
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
			ReturnErrorResponse(err, w)
			return
		}

		validationErr := ValidateRequest(requestProject)
		if validationErr != nil {
			ReturnErrorResponse(validationErr, w)
			return
		}

		projectId, err := createProject(
			database,
			requestProject.Name,
			requestProject.Type,
			requestProject.Client,
		)
		if err != nil {
			ReturnErrorResponse(err, w)
			return
		}

		response, err := getCreateProjectResponseBody(projectId)
		if err != nil {
			ReturnErrorResponse(err, w)
			return
		}

		w.Write(response)
	}
}

func CreateGetProjectsHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		projects, err := getProjects(database)
		if err != nil {
			ReturnErrorResponse(err, w)
			return
		}
		response, err := json.Marshal(projects)
		if err != nil {
			log.WithField("error", err.Error()).Error("Error marshaling the response")
			ReturnErrorResponse(&ErrorResponse{
				ErrorMessage: "Error mashaling the response",
				ErrorCode:    500,
			}, w)
			return
		}
		w.Write(response)
	}
}
