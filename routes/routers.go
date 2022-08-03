package routes

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(string) http.HandlerFunc
}

type Routes []Route

func NewRouter(config EnvConfiguration) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc("")
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"Healthiness",
		strings.ToUpper("Get"),
		"/-/health",
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
