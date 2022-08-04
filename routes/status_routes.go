package routes

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StatusResponse struct {
	Status  string `json:"status"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func CreateHealthinessHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		response := StatusResponse{
			Status:  "OK",
			Name:    "issue-service",
			Version: "1.0.0",
		}

		byteReponse, err := json.Marshal(response)

		if err != nil {
			log.WithField("error", err.Error()).Fatal("Error in healthiness probe")
		}
		w.Write(byteReponse)
	}
}

func CreateReadinessHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		response := StatusResponse{
			Status:  "OK",
			Name:    "issue-service",
			Version: "1.0.0",
		}

		byteReponse, err := json.Marshal(response)

		if err != nil {
			log.WithField("error", err.Error()).Fatal("Error in readiness probe")
		}

		w.Write(byteReponse)
	}
}
