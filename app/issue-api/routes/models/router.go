package models

import (
	"net/http"

	"gorm.io/gorm"
)

type ErrorResponse struct {
	ErrorMessage string
	ErrorCode    int
}

func (err ErrorResponse) Error() string {
	return err.ErrorMessage
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*gorm.DB) http.HandlerFunc
}

type Routes []Route

type Router interface {
	// Routes returns the list of routes to add to the docker server.
	Routes() Routes
}
