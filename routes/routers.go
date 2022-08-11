package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gorm.io/gorm"
)

const DUPLICATE_KEY_ERROR = "duplicate key value violates unique constraint"

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

func NewRouter(config EnvConfiguration) *negroni.Negroni {
	router := mux.NewRouter().StrictSlash(true)
	nRouter := negroni.New(negroni.NewRecovery())

	db, err := ConnectDatabase(config)

	if err != nil {
		return &negroni.Negroni{}
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc(db)
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	nRouter.UseHandler(router)
	return nRouter
}

func IsDuplicateKeyError(databaseError error) bool {
	return strings.Contains(databaseError.Error(), DUPLICATE_KEY_ERROR)
}

func ValidateRequest(inputRequest interface{}) error {
	var validationError = ""
	validate := validator.New()
	err := validate.Struct(inputRequest)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Validation error, field: %s, tag: %s", err.Namespace(), err.Tag())
			if validationError == "" {
				validationError = fmt.Sprintf("%s%s", validationError, errorMessage)
			} else {
				validationError = fmt.Sprintf("%s\n%s", validationError, errorMessage)
			}
		}
		return ErrorResponse{
			ErrorMessage: validationError,
			ErrorCode:    400,
		}
	}
	return nil
}

var routes = Routes{
	Route{
		"Healthiness",
		strings.ToUpper("Get"),
		"/-/healthz",
		CreateHealthinessHandler,
	},

	Route{
		"Healthiness",
		strings.ToUpper("Get"),
		"/-/ready",
		CreateReadinessHandler,
	},

	Route{
		"AddIssue",
		strings.ToUpper("Post"),
		"/v1/projects/{projectId}/sprints/{sprintId}/issues",
		CreateAddIssueHandler,
	},

	Route{
		"GetIssues",
		strings.ToUpper("Get"),
		"/v1/projects/{projectId}/sprints/{sprintId}/issues",
		CreateGetIssuesHandler,
	},

	Route{
		"GetIssuesById",
		strings.ToUpper("Get"),
		"/v1/projects/{projectId}/sprints/{sprintId}/issues/{issueId}",
		CreateGetIssuesByIdHandler,
	},

	Route{
		"PatchIssuesById",
		strings.ToUpper("Patch"),
		"/v1/projects/{projectId}/sprints/{sprintId}/issues/{issueId}",
		CreatePatchIssuesByIdHandler,
	},

	Route{
		"AddProject",
		strings.ToUpper("Post"),
		"/v1/projects",
		CreateAddProjectHandler,
	},

	Route{
		"AddSprint",
		strings.ToUpper("Post"),
		"/v1/projects/{projectId}/sprints",
		CreateAddSprintHandler,
	},

	Route{
		"PatchSprint",
		strings.ToUpper("Patch"),
		"/v1/projects/{projectId}/sprints/{sprintId}",
		CreatePatchSprintHandler,
	},
}
