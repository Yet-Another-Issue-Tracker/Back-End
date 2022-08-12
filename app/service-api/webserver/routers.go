package webserver

import (
	"issue-service/app/service-api/cfg"
	"issue-service/app/service-api/routes"
	"issue-service/app/service-api/routes/makes/models"
	"issue-service/internal"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func NewRouter(config cfg.EnvConfiguration) *negroni.Negroni {
	router := mux.NewRouter().StrictSlash(true)
	nRouter := negroni.New(negroni.NewRecovery())

	// todo receive this from args
	db, err := internal.ConnectDatabase(config)

	if err != nil {
		return &negroni.Negroni{}
	}

	for _, route := range routesToRegister {
		var handler http.Handler
		handler = route.HandlerFunc(db)
		handler = internal.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	nRouter.UseHandler(router)
	return nRouter
}

var routesToRegister = models.Routes{
	models.Route{
		Name:        "Healthiness",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/-/healthz",
		HandlerFunc: routes.CreateHealthinessHandler,
	},

	models.Route{
		Name:        "Healthiness",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/-/ready",
		HandlerFunc: routes.CreateReadinessHandler,
	},

	models.Route{
		Name:        "AddIssue",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues",
		HandlerFunc: routes.CreateAddIssueHandler,
	},

	models.Route{
		Name:        "GetIssues",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues",
		HandlerFunc: routes.CreateGetIssuesHandler,
	},

	models.Route{
		Name:        "GetIssuesById",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues/{issueId}",
		HandlerFunc: routes.CreateGetIssuesByIdHandler,
	},

	models.Route{
		Name:        "PatchIssuesById",
		Method:      strings.ToUpper("Patch"),
		Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues/{issueId}",
		HandlerFunc: routes.CreatePatchIssuesByIdHandler,
	},

	models.Route{
		Name:        "AddProject",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/v1/projects",
		HandlerFunc: routes.CreateAddProjectHandler,
	},

	models.Route{
		Name:        "GetProjects",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/v1/projects",
		HandlerFunc: routes.CreateGetProjectsHandler,
	},

	models.Route{
		Name:        "AddSprint",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/v1/projects/{projectId}/sprints",
		HandlerFunc: routes.CreateAddSprintHandler,
	},

	models.Route{
		Name:        "PatchSprint",
		Method:      strings.ToUpper("Patch"),
		Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}",
		HandlerFunc: routes.CreatePatchSprintHandler,
	},
}
