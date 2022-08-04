package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

var (
	errGeneric = errors.New("internal Server Error")
	// errBadRequest = errors.New("bad Request")
)

func createProject(database *gorm.DB, projectName string, projectType string, projectClient string) ([]byte, error) {
	project := Project{Name: projectName, Client: projectClient, Type: projectType}

	result := database.Create(&project)

	if result.Error != nil {
		log.Fatalf("Error creating new project %s", result.Error.Error())
		return nil, result.Error
	}
	response := CreateProjectResponse{Id: fmt.Sprint(project.ID)}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
func CreateAddProjectHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		projectId, err := createProject(database, "", "", "")

		if err != nil {
			log.Fatalln("failed response unmarshalling")
			http.Error(w, errGeneric.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(projectId)
	}
}
