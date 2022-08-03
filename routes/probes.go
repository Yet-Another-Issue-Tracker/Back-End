package routes

import (
	"net/http"

	"gorm.io/gorm"
)

func CreateHealthinessHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		response := []byte("ok")
		w.Write(response)
	}
}

func CreateReadinessHandler(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		response := []byte("ok")
		w.Write(response)
	}
}
