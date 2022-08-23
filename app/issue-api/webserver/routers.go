package webserver

import (
	"issue-service/app/issue-api/cfg"
	"issue-service/app/issue-api/routes"
	"issue-service/app/issue-api/routes/issue"
	"issue-service/app/issue-api/routes/models"
	"issue-service/app/issue-api/routes/project"
	"issue-service/app/issue-api/routes/sprint"
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

	routesToRegister = append(routesToRegister, project.NewRouter(db).Routes()...)
	routesToRegister = append(routesToRegister, sprint.NewRouter(db).Routes()...)
	routesToRegister = append(routesToRegister, issue.NewRouter(db).Routes()...)
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
		Name:        "GetIssuesById",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/v1/projects/{projectId}/sprints/{sprintId}/issues/{issueId}",
		HandlerFunc: routes.CreateGetIssuesByIdHandler,
	},
}
