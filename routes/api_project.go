package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	errGeneric = errors.New("internal Server Error")
	// errBadRequest = errors.New("bad Request")
)

func createProject() ([]byte, error) {
	response := CreateProjectResponse{
		Id: "123",
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
func CreateAddProjectHandler(databaseConnection string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		projectId, err := createProject()

		if err != nil {
			log.Fatalln("failed response unmarshalling")
			http.Error(w, errGeneric.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(projectId)
	}
}
