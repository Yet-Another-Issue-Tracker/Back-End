package routes

import (
	"net/http"
)

func CreateHealthinessHandler(connectionString string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		response := []byte("ok")
		w.Write(response)
	}
}

func CreateReadinessHandler(connectionString string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		response := []byte("ok")
		w.Write(response)
	}
}
